package seed

import (
	"log"
	"math/rand"
	"time"

	"backend/internal/app/entity"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	DefaultPrefixID uint = 3 // default (Admin / Advisor)
	DefaultMajorID  uint = 1 // fallback Major
)

func SeedUsers(db *gorm.DB) {
	// แนะนำ: rand.Seed ควรเรียกครั้งเดียวตอน start app แต่ใส่ไว้ที่นี่ก็ไม่ผิดครับ
	rand.Seed(time.Now().UnixNano())

	// ------------------------------
	// โหลด IT Department สำหรับ Admin
	// ------------------------------
	var itDept entity.Department
	if err := db.Where("department = ?", "IT Department").First(&itDept).Error; err != nil {
		log.Println("ไม่พบ IT Department (ต้อง seed department ก่อน)")
		return
	}

	// ------------------------------
	// โหลด Major สำหรับใส่ MajorID
	// ------------------------------
	var majors []entity.Major
	if err := db.Find(&majors).Error; err != nil || len(majors) == 0 {
		log.Println("ไม่พบ Major ใน DB")
		return
	}

	validMajorIDs := make([]uint, 0)
	for _, m := range majors {
		if m.Major != "ยังไม่สังกัดสาขา" {
			validMajorIDs = append(validMajorIDs, m.ID)
		}
	}
	if len(validMajorIDs) == 0 {
		validMajorIDs = append(validMajorIDs, DefaultMajorID)
	}

	// ------------------------------
	// สร้าง Admin
	// ------------------------------
	var roleAdmin entity.Role
	if err := db.Where("role = ?", "Admin").First(&roleAdmin).Error; err != nil {
		log.Println("ไม่พบ role Admin")
		return
	}

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("admin"), 14)
	admin := entity.User{
		SutId:        "A6608019",
		FirstName:    "Admin",
		LastName:     "System",
		Email:        "admin@example.com",
		PasswordHash: string(hashedPass),
		RoleID:       roleAdmin.ID,
		PrefixID:     DefaultPrefixID,
		MajorID:      validMajorIDs[0],
		DepartmentID: &itDept.ID,
		Active:       true,
	}

	// แก้ไข: ใช้ FirstOrCreate เพื่อป้องกัน Error Duplicate Key
	if err := db.Where("sut_id = ?", admin.SutId).FirstOrCreate(&admin).Error; err != nil {
		log.Fatalf("จัดการข้อมูล Admin ล้มเหลว: %v", err)
	}
	log.Println("Seed Admin เรียบร้อย (หรือมีอยู่แล้ว):", admin.SutId)

	// ------------------------------
	// สร้าง Advisor
	// ------------------------------
	var roleAdvisor entity.Role
	if err := db.Where("role = ?", "Advisor").First(&roleAdvisor).Error; err != nil {
		log.Println("ไม่พบ role Advisor")
		return
	}
	advisorPrefixes := []uint{4, 5, 6, 7, 8, 9, 10}

	// สุ่ม Prefix ใหม่ทุกครั้งที่รันอาจจะไม่ดีถ้า user มีอยู่แล้ว แต่นี่เป็น Seed ยอมรับได้
	advisorPrefixID := advisorPrefixes[rand.Intn(len(advisorPrefixes))]
	hashedAdvisorPass, _ := bcrypt.GenerateFromPassword([]byte("advisor123"), 14)
	
	advisor := entity.User{
		SutId:        "T6608019",
		FirstName:    "Advisor",
		LastName:     "User",
		Email:        "advisor@example.com",
		PasswordHash: string(hashedAdvisorPass),
		RoleID:       roleAdvisor.ID,
		PrefixID:     advisorPrefixID,
		MajorID:      validMajorIDs[0],
		ManagerID:    &admin.ID,
		Active:       true,
	}

	// แก้ไข: ใช้ FirstOrCreate
	if err := db.Where("sut_id = ?", advisor.SutId).FirstOrCreate(&advisor).Error; err != nil {
		log.Fatalf("จัดการข้อมูล Advisor ล้มเหลว: %v", err)
	}
	log.Println("Seed Advisor เรียบร้อย (หรือมีอยู่แล้ว):", advisor.SutId)

	// ------------------------------
	// สร้าง AdvisorProfile
	// ------------------------------
	advisorProfile := entity.AdvisorProfile{
		UserID:       advisor.ID, // advisor.ID จะมีค่าเสมอหลังจากผ่าน FirstOrCreate ข้างบน
		OfficeRoom:   "A101",
		Specialties:  "คอมพิวเตอร์",
	}

	// แก้ไข: เช็คว่ามี Profile ของ UserID นี้หรือยัง
	if err := db.Where("user_id = ?", advisor.ID).FirstOrCreate(&advisorProfile).Error; err != nil {
		log.Fatalf("สร้าง AdvisorProfile ล้มเหลว: %v", err)
	}
	log.Println("ตรวจสอบ AdvisorProfile เรียบร้อย")

	// ------------------------------
	// Seed Students
	// ------------------------------
	type seedStudent struct {
		SutId     string
		FirstName string
		LastName  string
		Email     string
	}

	students := []seedStudent{
		{"B6608019", "เนตรนภัทร", "ชำนินอก", "Netnaphat.cha@gmail.com"},
		{"B6608020", "ชยุต", "นนทการ", "chayut8020@example.com"},
		{"B6608021", "กัญญารัตน์", "พรมชัย", "kanyarat8021@example.com"},
		{"B6608022", "ธนภัทร", "เรืองกิจ", "thanaphat8022@example.com"},
		{"B6608023", "ณัฐชยา", "ศรีสมบัติ", "natchaya8023@example.com"},
		{"B6608024", "ภูวเดช", "กิตติวงศ์", "phuwadet8024@example.com"},
		{"B6608025", "พิมพ์อร", "จิตรานนท์", "pimorn8025@example.com"},
		{"B6608026", "อธิวัฒน์", "ทรงพล", "athiwat8026@example.com"},
		{"B6608027", "ชนิกานต์", "พูนศรี", "chanikan8027@example.com"},
		{"B6608028", "วชิรวิทย์", "สุนทรกุล", "wachirawit8028@example.com"},
		{"B6608029", "กุลิสรา", "มณีกาญจน์", "kulissara8029@example.com"},
	}

	var roleStudent entity.Role
	if err := db.Where("role = ?", "Student").First(&roleStudent).Error; err != nil {
		log.Println("ไม่พบ role Student")
		return
	}

	studentPrefixes := []uint{1, 3} // นาย, นางสาว

	for _, s := range students {
		prefixID := studentPrefixes[rand.Intn(len(studentPrefixes))]
		majorID := validMajorIDs[rand.Intn(len(validMajorIDs))]

		hashed, _ := bcrypt.GenerateFromPassword([]byte("student"), 14)
		user := entity.User{
			SutId:        s.SutId,
			FirstName:    s.FirstName,
			LastName:     s.LastName,
			Email:        s.Email,
			PasswordHash: string(hashed),
			RoleID:       roleStudent.ID,
			PrefixID:     prefixID,
			MajorID:      majorID,
			ManagerID:    &admin.ID, // assign admin เป็นผู้ดูแล
			Active:       true,
		}
		
		// แก้ไข: ใช้ FirstOrCreate โดยเช็คจาก sut_id
		if err := db.Where("sut_id = ?", s.SutId).FirstOrCreate(&user).Error; err != nil {
			log.Printf("จัดการข้อมูล Student %s ล้มเหลว: %v", s.SutId, err)
		} else {
			log.Println("Seed Student เรียบร้อย (หรือมีอยู่แล้ว):", s.SutId)
		}
	}

	log.Println("Seed Users เสร็จสมบูรณ์ ✅")
}