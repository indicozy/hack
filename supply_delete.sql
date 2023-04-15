DELIMITER $$

CREATE PROCEDURE UpdateSoldAmount(IN need_amount FLOAT, IN sale_price FLOAT)
BEGIN
  DECLARE finished INTEGER DEFAULT 0;
  DECLARE curr_barcode VARCHAR(255);
  DECLARE curr_quantity INTEGER;
  DECLARE curr_sold_amount FLOAT;
  DECLARE delta FLOAT;
  DECLARE cur_price FLOAT;
  DECLARE final_sum FLOAT DEFAULT 0;
  DECLARE need_amount_original FLOAT DEFAULT need_amount;

  -- Cursor to fetch rows that meet the criteria
  DECLARE cur CURSOR FOR
    SELECT barcode, quantity, sold_amount, price
    FROM supply
    WHERE sold_amount < quantity
    ORDER BY datetime ASC;

  -- Handler for 'not found'
  DECLARE CONTINUE HANDLER FOR NOT FOUND SET finished = 1;

  OPEN cur;

  update_loop: LOOP
    FETCH cur INTO curr_barcode, curr_quantity, curr_sold_amount, cur_price;

    IF finished = 1 THEN
      LEAVE update_loop;
    END IF;

    SET delta = curr_quantity - curr_sold_amount;
    IF delta <= need_amount THEN
      SET need_amount = need_amount - delta;
      SET curr_sold_amount = curr_quantity;
      SET final_sum = final_sum + (delta * cur_price);
    ELSE
      SET curr_sold_amount = curr_sold_amount + need_amount;
      SET final_sum = final_sum + (need_amount * cur_price);
      SET need_amount = 0;
      UPDATE supply
      SET sold_amount = curr_sold_amount
      WHERE barcode = curr_barcode;
      LEAVE update_loop;
    END IF;

    UPDATE supply
    SET sold_amount = curr_sold_amount
    WHERE barcode = curr_barcode;

  END LOOP update_loop;

  CLOSE cur;

  SELECT final_sum;
  
  INSERT INTO sale (barcode, quantity, datetime, price, margin)
  VALUES ('ABC123', need_amount_original, NOW(), sale_price, ROUND(need_amount_original * sale_price - final_sum, 2));
  
  SET @profit = ROUND(need_amount_original * sale_price - final_sum, 2);

  IF (SELECT COUNT(*) FROM daily_profit_table WHERE date = CURDATE()) > 0 THEN
    UPDATE daily_profit_table SET profit = profit + @profit WHERE date = CURDATE();
  ELSE
    INSERT INTO daily_profit_table (date, profit) VALUES (CURDATE(), @profit);
  END IF;
  
  
  SET @profit = ROUND(need_amount_original * sale_price - final_sum, 2);
  SET @year = YEAR(CURDATE());
  SET @month = MONTH(CURDATE());
  
  IF (SELECT COUNT(*) FROM monthly_profit_table WHERE year = @year AND month = @month) > 0 THEN
      UPDATE monthly_profit_table SET profit = profit + @profit WHERE year = @year AND month = @month;
  ELSE
      INSERT INTO monthly_profit_table (year, month, profit) VALUES (@year, @month, @profit);
  END IF;


END $$

DELIMITER ;


CALL UpdateSoldAmount(100, 10); -- python variables
