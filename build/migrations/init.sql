CREATE TABLE white_list(
    id     serial  not null,
    subnet cidr not null
);
CREATE UNIQUE INDEX white_list_subnet_index ON white_list (subnet);
CREATE TABLE black_list(
    id     serial  not null,
    subnet cidr not null
);
CREATE UNIQUE INDEX black_list_subnet_index ON black_list (subnet);