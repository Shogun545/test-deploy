package seed

import (
	"backend/internal/app/entity"
	"errors"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

func SeedAppointments(db *gorm.DB) error {
	log.Println("🌱 Seeding appointments...")

	// 1) หา role ID ของ ADVISOR/STUDENT โดยไม่ hardcode + auto-detect column
	advisorRoleID, err := findRoleIDAuto(db, []string{"ADVISOR", "Advisor", "advisor"})
	if err != nil {
		log.Println("❌ Advisor role not found, skip appointment seed")
		return nil
	}

	studentRoleID, err := findRoleIDAuto(db, []string{"STUDENT", "Student", "student"})
	if err != nil {
		log.Println("❌ Student role not found, skip appointment seed")
		return nil
	}

	// 2) หา user advisor/student
	var advisor entity.User
	if err := db.Where("role_id = ?", advisorRoleID).First(&advisor).Error; err != nil {
		log.Println("❌ Advisor not found, skip appointment seed")
		return nil
	}

	var student entity.User
	if err := db.Where("role_id = ?", studentRoleID).First(&student).Error; err != nil {
		log.Println("❌ Student not found, skip appointment seed")
		return nil
	}

	// 3) หา topic/category อย่างน้อย 1 อัน
	var topic entity.AppointmentTopic
	if err := db.First(&topic).Error; err != nil {
		log.Println("❌ Topic not found, skip appointment seed")
		return nil
	}

	var category entity.AppointmentCategory
	if err := db.First(&category).Error; err != nil {
		log.Println("❌ Category not found, skip appointment seed")
		return nil
	}

	// 4) ลบ appointment เดิม (ถ้ามี) เพื่อให้ seed สร้างใหม่ทุกครั้ง
	if err := db.Where(
		"advisor_user_id = ? AND student_user_id = ? AND topic_id = ?",
		advisor.ID, student.ID, topic.ID,
	).Delete(&entity.Appointment{}).Error; err != nil {
		log.Println("❌ Failed to delete existing appointment:", err)
		return err
	}

	appointment := entity.Appointment{
		Description:         "ขอปรึกษาเรื่อง Project จบ",
		AdvisorUserID:       advisor.ID,
		StudentUserID:       student.ID,
		TopicID:             topic.ID,
		CategoryID:          category.ID,
		AppointmentStatusID: entity.StatusPendingID,
	}

	if err := db.Create(&appointment).Error; err != nil {
		log.Println("❌ Failed to create appointment:", err)
		return err
	}

	log.Println("✅ Appointment seed completed")
	return nil
}

// -------------------- helpers --------------------

// findRoleIDAuto (quiet):
// 1) ตรวจคอลัมน์ของตาราง roles ผ่าน information_schema
// 2) เลือกคอลัมน์ที่น่าจะเป็นชื่อ role (เช่น name, role, role_name, title, code)
// 3) query หา role ด้วย candidates โดยใช้ Find() (ไม่เด้ง record not found)
func findRoleIDAuto(db *gorm.DB, candidates []string) (uint, error) {
	roleTable := db.NamingStrategy.TableName("roles") // ปกติจะเป็น "roles"

	cols, err := getTableColumns(db, roleTable)
	if err != nil {
		return 0, err
	}

	// คอลัมน์ที่มักจะเก็บชื่อ role
	preferred := []string{"name", "role", "role_name", "rolename", "title", "code", "type"}

	// หา column ตัวแรกที่มีจริงในตาราง
	col := pickFirstExistingColumn(cols, preferred)
	if col == "" {
		return 0, fmt.Errorf("no role-name column found in table %s (have: %v)", roleTable, cols)
	}

	// ลองหา role ตาม candidates (ใช้ Find เพื่อไม่ให้ log record not found)
	for _, c := range candidates {
		var role entity.Role
		if err := db.Where(fmt.Sprintf("%s = ?", col), c).Limit(1).Find(&role).Error; err != nil {
			// error จริงเท่านั้น
			return 0, err
		}
		if role.ID != 0 {
			return role.ID, nil
		}
	}

	// ไม่เจอจริง ๆ
	return 0, gorm.ErrRecordNotFound
}

func getTableColumns(db *gorm.DB, table string) ([]string, error) {
	var cols []string
	rows, err := db.Raw(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = ?
		ORDER BY ordinal_position
	`, table).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		cols = append(cols, c)
	}
	if len(cols) == 0 {
		return nil, errors.New("roles table columns not found (table may not exist or schema != public)")
	}
	return cols, nil
}

func pickFirstExistingColumn(have []string, preferred []string) string {
	set := map[string]bool{}
	for _, h := range have {
		set[strings.ToLower(h)] = true
	}
	for _, p := range preferred {
		if set[strings.ToLower(p)] {
			return p
		}
	}
	return ""
}
