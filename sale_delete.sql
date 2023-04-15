-- sale_id
UPDATE supply
  SET quantity = quantity - (SELECT quantity FROM supply WHERE id = <sale_id>) -- idk
  WHERE id IN (SELECT sale_id FROM supply_sale WHERE sale_id = <sale_id>);

DELETE FROM supply_sale WHERE sale_id=<sale_id>

DELETE FROM sale
  WHERE id = <sale_id>;