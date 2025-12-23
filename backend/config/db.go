package config

import (
    "backend/config/seed"
    "backend/internal/app/entity"
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var db *gorm.DB

func DB() *gorm.DB {
    return db
}

func ConnectDB() {
    if os.Getenv("GIN_MODE") != "release" {
        if err := godotenv.Load(); err != nil {
            log.Println("Warning: .env file not found, using environment variables only")
        }
    }

    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
        host, user, password, dbname, port,
    )

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }

    db = database
    log.Println("Connected to database")
}

func SetupDatabase() {
    // 1. ปิดการตรวจสอบ Foreign Key ชั่วคราว 
    //    PostgreSQL จะสร้างตารางได้ทั้งหมด แม้จะมีปัญหาลำดับการอ้างอิงกัน
    db.Exec("SET session_replication_role = 'replica';")
    
    // 2. Migrate ตารางทั้งหมดในครั้งเดียว
    //    เนื่องจากปิด Foreign Key แล้ว จึงไม่มีปัญหาลำดับการเรียก
    if err := db.AutoMigrate(
        &entity.User{},
        &entity.AcademicCalendar{},
        &entity.AdvisorNonAvailabillity{},
        
        &entity.TimeNonAvailabillity{},
        &entity.AdvisorProfile{},
        &entity.StudentProfile{},
        &entity.StudentAcademicRecord{},
        &entity.Appointment{},
        &entity.AdvisorLog{},
        &entity.ProgressReport{},
        &entity.ReportFeedback{},
        &entity.AppointmentReminder{},
        &entity.AppointmentState{},
        &entity.AppointmentStatus{},
        &entity.ApprovalAction{},
        &entity.Notification{},
        &entity.NotificationsHistory{},
        &entity.StatusHistory{},
        &entity.FAQ{},
        &entity.Report{}, 
        &entity.ReportImage{}, // รวม ReportImage เข้ามาใน Batch นี้ได้เลย
        &entity.ReportStatus{},
        &entity.ReportTopic{},
    ); err != nil {
        log.Fatalf("failed to migrate schema: %v", err)
    }

    // 3. เปิดการตรวจสอบ Foreign Key กลับมา
    db.Exec("SET session_replication_role = 'origin';")

    // 4. ส่วน Seed ข้อมูล (รันหลังจากตารางทั้งหมดถูกสร้าง)
    seed.SeedPrefix(db)
    seed.SeedMajor(db)
    seed.SeedRoleData(db)
	seed.SeedDepartmentData(db)
    seed.SeedUsers(db)
    seed.SeedStudentProfile(db)
    seed.SeedAppointmentStatus(db)
    seed.SeedAdvisorProfile(db) 
    seed.SeedApprovalActions(db)
    seed.SeedHolidays(db)
    seed.SeedAppointments(db)
    seed.SeedAppointmentCategories(db)
    seed.SeedAppointmentTopics(db)
    seed.SeedReportStatus(db)
    seed.SeedReportTopic(db)
    seed.SeedReport(db)

    log.Println("Database migration and seeding complete! Server Ready.")
}