-- Ordering requires a stock record per menu item, but items created before
-- menu creation added one have none and can never be ordered or restocked.
-- Give each of them a record at zero stock for the admin to restock.
INSERT INTO inventory_items (menu_item_id, name, current_stock)
SELECT m.id, m.name, 0
FROM menu_items m
WHERE NOT EXISTS (
    SELECT 1 FROM inventory_items i WHERE i.menu_item_id = m.id
);
