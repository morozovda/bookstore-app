CREATE TABLE "book" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title VARCHAR(200) NOT NULL,
    author VARCHAR(200) NOT NULL,
    price INT NOT NULL,
    amount INT NOT NULL
);

INSERT INTO "book" ("title", "author", "price", "amount") VALUES
('Alice''s Adventures in Wonderland', 'Lewis Carroll', 400, 15),
('The Selfish Gene', 'Richard Dawkins', 600, 7),
('Dead Souls', 'Nikolai Gogol', 400, 11),
('War and Peace', 'Leo Tolstoy', 700, 5),
('The Brothers Karamazov', 'Fyodor Dostoevsky', 400, 19),
('Sword of Destiny', 'Andrzej Sapkowski', 300, 16),
('Death Note', 'Tsugumi Ohba', 500, 10),
('Berserk', 'Kentaro Miura', 600, 5),
('Tokyo Ghoul', 'Sui Ishida', 300, 6);

CREATE TABLE "customer" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(200) NOT NULL,
    email VARCHAR(200) NOT NULL,
    passwd VARCHAR(60) NOT NULL,
    balance INT DEFAULT 0
);

CREATE TABLE deal (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    book_id UUID NOT NULL,
    order_amount INT NOT NULL,
    customer_id UUID NOT NULL,
    CONSTRAINT fk_book FOREIGN KEY(book_id) REFERENCES book(id),
    CONSTRAINT fk_customer FOREIGN KEY(customer_id) REFERENCES customer(id)
);