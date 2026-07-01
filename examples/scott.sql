create table dept (deptno integer primary key, dname varchar, loc varchar);
create table emp (empno integer primary key, ename varchar, job varchar, mgr integer, sal numeric, comm numeric, deptno integer);
create table bonus (id integer primary key, ename varchar, job varchar, sal numeric, comm numeric);
create table salgrade (grade integer primary key, losal numeric, hisal numeric);

insert into dept values (10, 'ACCOUNTING', 'NEW YORK');
insert into dept values (20, 'RESEARCH', 'DALLAS');
insert into dept values (30, 'SALES', 'CHICAGO');
insert into dept values (40, 'OPERATIONS', 'BOSTON');

insert into emp values (7839, 'KING', 'PRESIDENT', 0, 5000, 0, 10);
insert into emp values (7698, 'BLAKE', 'MANAGER', 7839, 2850, 0, 30);
insert into emp values (7782, 'CLARK', 'MANAGER', 7839, 2450, 0, 10);
insert into emp values (7566, 'JONES', 'MANAGER', 7839, 2975, 0, 20);
insert into emp values (7788, 'SCOTT', 'ANALYST', 7566, 3000, 0, 20);
insert into emp values (7902, 'FORD', 'ANALYST', 7566, 3000, 0, 20);
insert into emp values (7369, 'SMITH', 'CLERK', 7902, 800, 0, 20);
insert into emp values (7499, 'ALLEN', 'SALESMAN', 7698, 1600, 300, 30);
insert into emp values (7521, 'WARD', 'SALESMAN', 7698, 1250, 500, 30);
insert into emp values (7654, 'MARTIN', 'SALESMAN', 7698, 1250, 1400, 30);
insert into emp values (7844, 'TURNER', 'SALESMAN', 7698, 1500, 0, 30);
insert into emp values (7876, 'ADAMS', 'CLERK', 7788, 1100, 0, 20);
insert into emp values (7900, 'JAMES', 'CLERK', 7698, 950, 0, 30);
insert into emp values (7934, 'MILLER', 'CLERK', 7782, 1300, 0, 10);

insert into salgrade values (1, 700, 1200);
insert into salgrade values (2, 1201, 1400);
insert into salgrade values (3, 1401, 2000);
insert into salgrade values (4, 2001, 3000);
insert into salgrade values (5, 3001, 9999);