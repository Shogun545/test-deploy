package seed

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"backend/internal/app/entity"
)

// ข้อมูลวันหยุด (แปลงจาก route.ts)
var holidays = map[string][]struct {
	Date      string
	LocalName string
	Name      string
}{
	"2025": {
		{"2025-01-01", "วันขึ้นปีใหม่", "New Year's Day"},
		{"2025-02-12", "วันมาฆบูชา", "Makha Bucha"},
		{"2025-04-06", "วันจักรี", "Chakri Memorial Day"},
		{"2025-04-13", "วันสงกรานต์", "Songkran Festival"},
		{"2025-04-14", "วันสงกรานต์", "Songkran Festival"},
		{"2025-04-15", "วันสงกรานต์", "Songkran Festival"},
		{"2025-05-04", "วันฉัตรมงคล", "Coronation Day"},
		{"2025-05-11", "วันวิสาขบูชา", "Visakha Bucha"},
		{"2025-06-03", "วันเฉลิมพระชนมพรรษาราชินี", "HM Queen Suthida's Birthday"},
		{"2025-07-10", "วันอาสาฬหบูชา", "Asalha Bucha"},
		{"2025-07-11", "วันเข้าพรรษา", "Buddhist Lent Day"},
		{"2025-07-28", "วันเฉลิมพระชนมพรรษา ร.10", "HM King Maha Vajiralongkorn's Birthday"},
		{"2025-08-12", "วันแม่แห่งชาติ", "Mother's Day"},
		{"2025-10-13", "วันคล้ายวันสวรรคต ร.9", "HM King Bhumibol Adulyadej Memorial Day"},
		{"2025-10-23", "วันปิยมหาราช", "Chulalongkorn Day"},
		{"2025-12-05", "วันพ่อแห่งชาติ", "Father's Day"},
		{"2025-12-10", "วันรัฐธรรมนูญ", "Constitution Day"},
		{"2025-12-31", "วันสิ้นปี", "New Year's Eve"},
	},
	"2026": {
		{"2026-01-01", "วันขึ้นปีใหม่", "New Year's Day"},
		{"2026-03-03", "วันมาฆบูชา", "Makha Bucha"},
		{"2026-04-06", "วันจักรี", "Chakri Memorial Day"},
		{"2026-04-13", "วันสงกรานต์", "Songkran Festival"},
		{"2026-04-14", "วันสงกรานต์", "Songkran Festival"},
		{"2026-04-15", "วันสงกรานต์", "Songkran Festival"},
		{"2026-05-04", "วันฉัตรมงคล", "Coronation Day"},
		{"2026-05-31", "วันวิสาขบูชา", "Visakha Bucha"},
		{"2026-06-03", "วันเฉลิมพระชนมพรรษาราชินี", "HM Queen Suthida's Birthday"},
		{"2026-07-28", "วันเฉลิมพระชนมพรรษา ร.10", "HM King Maha Vajiralongkorn's Birthday"},
		{"2026-07-29", "วันอาสาฬหบูชา", "Asalha Bucha"},
		{"2026-07-30", "วันเข้าพรรษา", "Buddhist Lent Day"},
		{"2026-08-12", "วันแม่แห่งชาติ", "Mother's Day"},
		{"2026-10-13", "วันคล้ายวันสวรรคต ร.9", "HM King Bhumibol Adulyadej Memorial Day"},
		{"2026-10-23", "วันปิยมหาราช", "Chulalongkorn Day"},
		{"2026-12-05", "วันพ่อแห่งชาติ", "Father's Day"},
		{"2026-12-10", "วันรัฐธรรมนูญ", "Constitution Day"},
		{"2026-12-31", "วันสิ้นปี", "New Year's Eve"},
	},
	"2027": {
		{"2027-01-01", "วันขึ้นปีใหม่", "New Year's Day"},
		{"2027-02-21", "วันมาฆบูชา", "Makha Bucha"},
		{"2027-04-06", "วันจักรี", "Chakri Memorial Day"},
		{"2027-04-13", "วันสงกรานต์", "Songkran Festival"},
		{"2027-04-14", "วันสงกรานต์", "Songkran Festival"},
		{"2027-04-15", "วันสงกรานต์", "Songkran Festival"},
		{"2027-05-04", "วันฉัตรมงคล", "Coronation Day"},
		{"2027-05-20", "วันวิสาขบูชา", "Visakha Bucha"},
		{"2027-06-03", "วันเฉลิมพระชนมพรรษาราชินี", "HM Queen Suthida's Birthday"},
		{"2027-07-18", "วันอาสาฬหบูชา", "Asalha Bucha"},
		{"2027-07-19", "วันเข้าพรรษา", "Buddhist Lent Day"},
		{"2027-07-28", "วันเฉลิมพระชนมพรรษา ร.10", "HM King Maha Vajiralongkorn's Birthday"},
		{"2027-08-12", "วันแม่แห่งชาติ", "Mother's Day"},
		{"2027-10-13", "วันคล้ายวันสวรรคต ร.9", "HM King Bhumibol Adulyadej Memorial Day"},
		{"2027-10-23", "วันปิยมหาราช", "Chulalongkorn Day"},
		{"2027-12-05", "วันพ่อแห่งชาติ", "Father's Day"},
		{"2027-12-10", "วันรัฐธรรมนูญ", "Constitution Day"},
		{"2027-12-31", "วันสิ้นปี", "New Year's Eve"},
	},
	"2028": {
		{"2028-01-01", "วันขึ้นปีใหม่", "New Year's Day"},
		{"2028-02-10", "วันมาฆบูชา", "Makha Bucha"},
		{"2028-04-06", "วันจักรี", "Chakri Memorial Day"},
		{"2028-04-13", "วันสงกรานต์", "Songkran Festival"},
		{"2028-04-14", "วันสงกรานต์", "Songkran Festival"},
		{"2028-04-15", "วันสงกรานต์", "Songkran Festival"},
		{"2028-05-04", "วันฉัตรมงคล", "Coronation Day"},
		{"2028-05-08", "วันวิสาขบูชา", "Visakha Bucha"},
		{"2028-06-03", "วันเฉลิมพระชนมพรรษาราชินี", "HM Queen Suthida's Birthday"},
		{"2028-07-06", "วันอาสาฬหบูชา", "Asalha Bucha"},
		{"2028-07-07", "วันเข้าพรรษา", "Buddhist Lent Day"},
		{"2028-07-28", "วันเฉลิมพระชนมพรรษา ร.10", "HM King Maha Vajiralongkorn's Birthday"},
		{"2028-08-12", "วันแม่แห่งชาติ", "Mother's Day"},
		{"2028-10-13", "วันคล้ายวันสวรรคต ร.9", "HM King Bhumibol Adulyadej Memorial Day"},
		{"2028-10-23", "วันปิยมหาราช", "Chulalongkorn Day"},
		{"2028-12-05", "วันพ่อแห่งชาติ", "Father's Day"},
		{"2028-12-10", "วันรัฐธรรมนูญ", "Constitution Day"},
		{"2028-12-31", "วันสิ้นปี", "New Year's Eve"},
	},
	"2029": {
		{"2029-01-01", "วันขึ้นปีใหม่", "New Year's Day"},
		{"2029-02-28", "วันมาฆบูชา", "Makha Bucha"},
		{"2029-04-06", "วันจักรี", "Chakri Memorial Day"},
		{"2029-04-13", "วันสงกรานต์", "Songkran Festival"},
		{"2029-04-14", "วันสงกรานต์", "Songkran Festival"},
		{"2029-04-15", "วันสงกรานต์", "Songkran Festival"},
		{"2029-05-04", "วันฉัตรมงคล", "Coronation Day"},
		{"2029-05-27", "วันวิสาขบูชา", "Visakha Bucha"},
		{"2029-06-03", "วันเฉลิมพระชนมพรรษาราชินี", "HM Queen Suthida's Birthday"},
		{"2029-07-25", "วันอาสาฬหบูชา", "Asalha Bucha"},
		{"2029-07-26", "วันเข้าพรรษา", "Buddhist Lent Day"},
		{"2029-07-28", "วันเฉลิมพระชนมพรรษา ร.10", "HM King Maha Vajiralongkorn's Birthday"},
		{"2029-08-12", "วันแม่แห่งชาติ", "Mother's Day"},
		{"2029-10-13", "วันคล้ายวันสวรรคต ร.9", "HM King Bhumibol Adulyadej Memorial Day"},
		{"2029-10-23", "วันปิยมหาราช", "Chulalongkorn Day"},
		{"2029-12-05", "วันพ่อแห่งชาติ", "Father's Day"},
		{"2029-12-10", "วันรัฐธรรมนูญ", "Constitution Day"},
		{"2029-12-31", "วันสิ้นปี", "New Year's Eve"},
	},
}

func SeedHolidays(db *gorm.DB) {
	fmt.Println("Start seeding holidays...")

	// ตรวจสอบว่ามี Admin อย่างน้อย 1 คนหรือไม่
	var admin entity.User
	// ใช้ Unscoped เพื่อหา Admin แม้ว่าจะถูก Soft Delete ไปแล้ว (กัน error)
	if err := db.Unscoped().First(&admin).Error; err != nil {
		fmt.Println("Warning: No user found. Using ID 1 as default.")
		admin.ID = 1
	}

	for _, yearEvents := range holidays {
		for _, event := range yearEvents {
			// 1. แปลง String date เป็น Time Object
			dateObj, err := time.Parse("2006-01-02", event.Date)
			if err != nil {
				continue
			}

			// 2. กำหนดเวลาเริ่ม-จบ (วันหยุดคือทั้งวัน 00:00 - 23:59)
			// ใช้ time.Local เพื่อให้ตรงกับ Timezone เครื่อง server
			startTime := time.Date(dateObj.Year(), dateObj.Month(), dateObj.Day(), 0, 0, 0, 0, time.Local)
			endTime := time.Date(dateObj.Year(), dateObj.Month(), dateObj.Day(), 23, 59, 59, 0, time.Local)

			// 3. ตรวจสอบว่ามีวันหยุดนี้อยู่แล้วหรือไม่ (เปลี่ยนเงื่อนไขจาก date เป็น start_date_time)
			var existCount int64
			db.Model(&entity.AcademicCalendar{}).
				Where("start_date_time = ? AND event_name = ?", startTime, event.LocalName).
				Count(&existCount)

			if existCount > 0 {
				continue
			}

			// 4. สร้าง AcademicCalendar (ตัด TimeCalendar ออก และใช้ฟิลด์ใหม่)
			academicCalendar := entity.AcademicCalendar{
				EventName:     event.LocalName,
				EventType:     "Public Holiday",
				// ✅ ใช้ 2 ฟิลด์ใหม่นี้แทน Date และ TimeCalendarID
				StartDateTime: startTime,
				EndDateTime:   endTime,
				
				IsHoliday:     true,
				AdminID:       admin.ID,
			}

			if err := db.Create(&academicCalendar).Error; err != nil {
				fmt.Printf("Failed to seed %s: %v\n", event.LocalName, err)
			} else {
				fmt.Printf("Seeded: %s (%s)\n", event.LocalName, event.Date)
			}
		}
	}
	fmt.Println("Holidays seeding completed.")
}