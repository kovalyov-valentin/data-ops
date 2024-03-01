CREATE TABLE IF NOT EXISTS logs(
  Id int,
  ProjectId int,
  Name VARCHAR(255),
  Description VARCHAR(255),
  Priority int,
  Removed bool,
  EventTime datetime
) ENGINE=MergeTree()
order by (Id);
insert into logs(Id, ProjectId, Name, Description, Priority, Removed, EventTime)
values (1,1,'1','',1,false,now());