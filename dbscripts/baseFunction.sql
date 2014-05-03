DROP FUNCTION if exists generateCategoryId;
DELIMITER $$
CREATE FUNCTION generateCategoryId (parentId VARCHAR(20)) RETURNS VARCHAR(50)
  begin
    declare tempId varchar(20);
    if parentId is null or parentId = '' or parentId = 0 then
      select max(id) + 1 into tempId from category;
      set tempId = ifnull(tempId, '101');
    else
      select cast(max(cast(id as unsigned)) + 1 as char(20)) into tempId from category where parent_id = parentId;
      if tempId is null then
        select concat(parentId, '001') into tempId;
      end if;
    end if;
    RETURN tempId;
  END$$
DELIMITER ;