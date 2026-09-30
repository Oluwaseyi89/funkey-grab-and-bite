-- Default menu categories the storefront is built around. Matching is
-- case-insensitive so an admin-created "noodles" is adopted, not duplicated.
WITH defaults (name, description, display_order) AS (
    VALUES
        ('Chips & Chicken', 'Crispy golden fries with our signature chicken', 1),
        ('Noodles', 'Delicious noodle dishes with fresh ingredients', 2),
        ('Shawarma', 'Authentic shawarma wraps and plates', 3),
        ('Drinks', 'Refreshing beverages', 4),
        ('Soup & Food Bowls', 'Hearty soups and nutritious bowls (Pre-order)', 5),
        ('Lunch Packs', 'Complete meal deals for lunch', 6)
),
updated AS (
    UPDATE menu_categories c
    SET description = d.description,
        display_order = d.display_order,
        is_active = true
    FROM defaults d
    WHERE LOWER(c.name) = LOWER(d.name)
)
INSERT INTO menu_categories (name, description, display_order, is_active)
SELECT d.name, d.description, d.display_order, true
FROM defaults d
WHERE NOT EXISTS (
      SELECT 1 FROM menu_categories c WHERE LOWER(c.name) = LOWER(d.name)
  )
ORDER BY d.display_order;
