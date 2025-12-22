# Table Schema

orders

- order_id
- user_id
- total_amount
- created_at

# Query

```sql
SELECT
    user_id,
    SUM(total_amount) AS total_spending
FROM orders
WHERE created_at >= CURRENT_TIMESTAMP - INTERVAL '30 days'
GROUP BY user_id;
```

# Indexing

```sql
CREATE INDEX idx_orders_created_at_user_id
ON orders (created_at, user_id);
```

Untuk optimalisasi query, dapat dibuat composite index antara created_at dan user_id. Hal itu dilakukan agar database dapat melakukan index range scan dan data yang telah diambil sudah berututan berdasarkan user_id. Sehingga ketika melakukan grouping by user_id tidak memakan resource dan waktu.
