-- Only remove seeded categories that no menu item still points at.
DELETE FROM menu_categories c
WHERE LOWER(c.name) IN (
        'chips & chicken',
        'noodles',
        'shawarma',
        'drinks',
        'soup & food bowls',
        'lunch packs'
    )
  AND NOT EXISTS (
      SELECT 1 FROM menu_items m WHERE m.category_id = c.id
  );
