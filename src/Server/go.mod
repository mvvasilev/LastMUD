module code.haedhutner.dev/mvv/LastMUD/Server

require (
    code.haedhutner.dev/mvv/LastMUD/CommandLib v0.0.0
    code.haedhutner.dev/mvv/LastMUD/CoreLib v0.0.0
)

replace (
    code.haedhutner.dev/mvv/LastMUD/CommandLib => ../CommandLib
    code.haedhutner.dev/mvv/LastMUD/CoreLib => ../CoreLib
)

go 1.24.4
