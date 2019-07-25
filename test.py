ProcInAirTemp:27 
ProcInFlo:74 
ProcInRH:22 
ProcOutAirTemp:29 
ReactInFlo:71 
ReactInTemp:131 
ReactOutTemp:98 
suid:1 
ts:06/24/2019 14:02:58.07




"INSERT INTO [Marposs].[dbo].[sgcdBryAir] 
( [ts] ,
[suid] ,
[ProcInFlo] ,
[ReactInFlo] ,
[SupAirRH] ,
[ProcInRH] ,
[ProcInAirTemp] ,
[ProcOutAirTemp] ,
[ReactInTemp] ,
[ReactOutTemp] )

VALUES(?,?,?,?,?,?,?,?)",

"[   {{sql_varchar, 80},  [emqx_plugin_2db:try_get_ts(Message, MessageMaps)] },
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"suid\">>, MessageMaps,\"int\")]},
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"ProcInFlo\">>, MessageMaps,\"int\")]},
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"ReactInFlo\">>, MessageMaps,\"int\")]},
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"SupAirRH\">>, MessageMaps,\"int\")]},
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"ProcInRH\">>, MessageMaps,\"int\")]},
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"ProcInAirTemp\">>, MessageMaps,\"int\")]},
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"ProcOutAirTemp\">>, MessageMaps,\"int\")]},
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"ReactInTemp\">>, MessageMaps,\"int\")]},
{sql_integer , [emqx_plugin_2db:try_get_val(<<\"ReactOutTemp\">>, MessageMaps,\"int\")]}
]\.",

SELECT [id]
      ,[ts]
      ,[suid]
      ,[ProcInFlo]
      ,[ReactInFlo]
      ,[SupAirRH]
      ,[ProcInRH]
      ,[ProcInAirTemp]
      ,[ProcOutAirTemp]
      ,[ReactInTemp]
      ,[ReactOutTemp]
  FROM [Marposs].[dbo].[sgcdBryAir] order by id desc