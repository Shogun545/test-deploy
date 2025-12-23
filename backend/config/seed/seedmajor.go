package seed

import (
    "log"

    "gorm.io/gorm"
    "backend/internal/app/entity"
)

func SeedMajor(db *gorm.DB) {

    majors := []entity.Major{
        {Major: "วิศวกรรมการผลิตอัตโนมัติและหุ่นยนต์"},
		{Major: "วิศวกรรมเกษตรและอาหาร"},
		{Major: "วิศวกรรมคอมพิวเตอร์"},
		{Major: "วิศวกรรมเคมี"},
		{Major: "วิศวกรรมเครื่องกล"},
		{Major: "วิศวกรรมปิโตรเลียมและเทคโนโลยีธรณี"},
		{Major: "วิศวกรรมไฟฟ้า"},
		{Major: "วิศวกรรมโทรคมนาคม"},
		{Major: "วิศวกรรมยานยนต์"},
		{Major: "วิศวกรรมโยธา"},
		{Major: "วิศวกรรมสิ่งแวดล้อม"},
		{Major: "วิศวกรรมอุตสาหการ"},
		{Major: "วิศวกรรมโลหการ"},
		{Major: "วิศวกรรมอิเล็กทรอนิกส์"},
		{Major: "วิศวกรรมขนส่งและโลจิสติกส์"},
		{Major: "วิศวกรรมเซรามิก"},
		{Major: "วิศวกรรมพอลิเมอร์"},
		{Major: "ยังไม่สังกัดสาขา"},
    }

    for _, m := range majors {
        var existing entity.Major
        err := db.Where("major = ?", m.Major).First(&existing).Error

        if err == gorm.ErrRecordNotFound {
            db.Create(&m)
        } else if err == nil {
            log.Printf("Major '%s' already exists, skipping...", m.Major)
        }
    }

    log.Println("Major seed complete.")
}
