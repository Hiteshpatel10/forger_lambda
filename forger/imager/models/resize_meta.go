package models

type ResizeMetaModel struct {
	Size    int
	DirName string
}

var AndroidResizeMetaList = []ResizeMetaModel{
	{Size: 48, DirName: "mipmap-mdpi"},
	{Size: 72, DirName: "mipmap-hdpi"},
	{Size: 96, DirName: "mipmap-xhdpi"},
	{Size: 144, DirName: "mipmap-xxhdpi"},
	{Size: 192, DirName: "mipmap-xxxhdpi"},
}

var IOSResizeMetaList = []ResizeMetaModel{
	{Size: 120, DirName: "AppIcon.appiconset"},
	{Size: 167, DirName: "AppIcon.appiconset"},
	{Size: 152, DirName: "AppIcon.appiconset"},
	{Size: 80, DirName: "AppIcon.appiconset"},
	{Size: 58, DirName: "AppIcon.appiconset"},
	{Size: 76, DirName: "AppIcon.appiconset"},
	{Size: 180, DirName: "AppIcon.appiconset"},
	{Size: 87, DirName: "AppIcon.appiconset"},
	{Size: 114, DirName: "AppIcon.appiconset"},
}
