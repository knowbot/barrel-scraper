CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS metadata (
      key VARCHAR PRIMARY KEY,
      last_updated VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS subcategories (
    id INTEGER PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE,
    url VARCHAR NOT NULL,
    category_id INTEGER NOT NULL,
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE IF NOT EXISTS provinces (
    id INTEGER PRIMARY KEY,
    name VARCHAR NOT NULL,
    code CHAR(2) NOT NULL UNIQUE,
    region VARCHAR NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_region ON provinces (region);

CREATE TABLE IF NOT EXISTS companies (
    id INTEGER PRIMARY KEY,
    name VARCHAR NOT NULL,
    streetAddr VARCHAR NOT NULL,
    cap VARCHAR NOT NULL,
    city VARCHAR NOT NULL,
    phone VARCHAR,
    fax VARCHAR,
    website VARCHAR,
    province_id INTEGER NOT NULL,
    subcategory_id INTEGER NOT NULL,
    CONSTRAINT fk_province FOREIGN KEY (province_id) REFERENCES provinces (id) CONSTRAINT fk_subcategory FOREIGN KEY (subcategory_id) REFERENCES subcategories (id)
);