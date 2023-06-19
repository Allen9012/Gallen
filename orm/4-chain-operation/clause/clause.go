/**
  @author: Allen
  @since: 2023/5/14
  @desc: //TODO
**/
package clause

// Support types for Clause
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)
