CREATE TABLE IF NOT EXISTS members (
    id VARCHAR(255) PRIMARY KEY,
    join_date DATE,
    date_of_birth DATE,
    city VARCHAR(255),
    no_of_child INT,
    eldest_kid_dob DATE,
    youngest_kid_dob DATE,
    password VARCHAR(255)
);