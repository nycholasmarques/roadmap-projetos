CREATE TABLE expense (
  id   SERIAL PRIMARY KEY,
  description text NOT NULL,
  amount  DECIMAL(10,2) NOT NULL,
  created_at date NOT NULL DEFAULT CURRENT_DATE
);
