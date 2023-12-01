module github.com/bkzy/micscript

go 1.16

replace (
	github.com/bkzy/micscript => ../micscript
	github.com/bkzy/micscript/engineauth/hardinfo => ../micscript/engineauth/hardinfo
)

require (
	github.com/StackExchange/wmi v1.2.0
	github.com/bkzy-wangjp/Author v1.0.0
	github.com/bkzy-wangjp/CRC16 v1.0.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/xuri/excelize/v2 v2.4.1
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1
)
