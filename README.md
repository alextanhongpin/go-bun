# go-bun


Sample dataset:

```sql
drop table authors;
create table if not exists authors (
	id int generated always as identity,
	name text not null,
	primary key (id)
);

insert into authors(name) values ('john'), ('jane');
table authors;

drop table if exists books;
create table if not exists books (
	id int generated always as identity,
	author_id int not null,
	title text not null,
	primary key(id),
	foreign key (author_id) references authors
);

insert into books (title, author_id) values
('Programming Go', 1),
('Programming Og', 2);
```
