-- supply_id
UPDATE sale
  SET margin = margin + (SELECT quantity * price FROM supply WHERE id = <supply_id>) -- idk
  WHERE id IN (SELECT sale_id FROM supply_sale WHERE supply_id = <supply_id>);

DELETE FROM supply_sale WHERE supply_id=<supply_id>;

DELETE FROM supply
  WHERE id = <supply_id>;
