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
  DECLARE @supply_id INTEGER;

  -- Cursor to fetch rows that meet the criteria
  DECLARE cur CURSOR FOR
    SELECT id, barcode, quantity, sold_amount, price
      FROM supply
      WHERE sold_amount < quantity
      ORDER BY datetime ASC;

  INSERT INTO sale (barcode, quantity, datetime, price, margin) -- тут потому что его ID нужен в дальнейшем
    VALUES ('ABC123', need_amount_original, NOW(), sale_price, ROUND(need_amount_original * sale_price - final_sum, 2));

  set @sale_id := SELECT LAST_INSERT_ID();

  -- Handler for 'not found'
  DECLARE CONTINUE HANDLER FOR NOT FOUND SET finished = 1;

  OPEN cur;

  update_loop: LOOP
    FETCH cur INTO @supply_id, curr_barcode, curr_quantity, curr_sold_amount, cur_price;

    IF finished = 1 THEN
      LEAVE update_loop;
    END IF;

    SET delta = curr_quantity - curr_sold_amount;
    IF delta <= need_amount THEN
      SET need_amount = need_amount - delta;
      SET curr_sold_amount = curr_quantity;
      SET final_sum = final_sum + (delta * cur_price);
    ELSE
      SET need_amount = 0;
      SET curr_sold_amount = curr_sold_amount + need_amount;
      SET final_sum = final_sum + (need_amount * cur_price);
      UPDATE supply
        SET sold_amount = curr_sold_amount
        WHERE barcode = curr_barcode;
      INSERT INTO supply_sale (supply_id, sale_id, supply_quantity)
        VALUES (@supply_id, @sale_id, curr_sold_amount);
      LEAVE update_loop;
    END IF;

    UPDATE supply
      SET sold_amount = curr_sold_amount
      WHERE barcode = curr_barcode;

    INSERT INTO supply_sale (supply_id, sale_id, supply_quantity)
      VALUES (@supply_id, @sale_id, curr_sold_amount);

  END LOOP update_loop;

  CLOSE cur;

  SELECT final_sum; -- returns final sum of sold amount

END $$

DELIMITER ;

CALL UpdateSoldAmount(100, 10); -- python variables
