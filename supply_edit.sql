-- new_quantity, supply_id

UPDATE sale 
  SET margin = margin + (SELECT quantity * price FROM supply WHERE id = supply_id)
  WHERE id IN (
    SELECT sale_id 
    FROM supply_sale 
    WHERE supply_id = supply_id
  );

UPDATE supply
  SET quantity = new_quantity
  WHERE id = supply_id;
