CREATE TABLE "book" (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    author VARCHAR(200) NOT NULL,
    price INT NOT NULL,
    amount INT NOT NULL
);

INSERT INTO "book" ("title", "author", "price", "amount") VALUES
('Alice''s Adventures in Wonderland', 'Lewis Carroll', 4, 15),
('The Selfish Gene', 'Richard Dawkins', 6, 7),
('Dead Souls', 'Nikolai Gogol', 4, 11),
('War and Peace', 'Leo Tolstoy', 7, 5),
('The Brothers Karamazov', 'Fyodor Dostoevsky', 4, 19),
('Sword of Destiny', 'Andrzej Sapkowski', 3, 16),
('Death Note', 'Tsugumi Ohba', 5, 10),
('Berserk', 'Kentaro Miura', 6, 5),
('Tokyo Ghoul', 'Sui Ishida', 3, 6);

CREATE TABLE "customer" (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    passwd VARCHAR(60) NOT NULL,
    balance INT DEFAULT 0
);

CREATE TABLE deal (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    order_amount INT NOT NULL,
    customer_id INT NOT NULL,
    CONSTRAINT fk_book FOREIGN KEY(book_id) REFERENCES book(id),
    CONSTRAINT fk_customer FOREIGN KEY(customer_id) REFERENCES customer(id)
);