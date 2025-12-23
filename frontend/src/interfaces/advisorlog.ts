export interface AdvisorLogInterface {
  ID?: number;

  /* ---------- Relation ---------- */
  AppointmentID?: number;

  /* ---------- Main Content ---------- */
  Title?: string;
  Body?: string;

  /* ---------- Status ---------- */
  Status?: AdvisorLogStatus;
  RequiresReport?: boolean;

  /* ---------- Files (comma-separated from backend) ---------- */
  FileName?: string; // "a.pdf,b.jpg"
  FilePath?: string; // "uploads/uuid1.pdf,uploads/uuid2.png"

  /* ---------- Timestamps ---------- */
  CreatedAt?: string;
  UpdatedAt?: string;
}

export type AdvisorLogStatus =
  | "Draft"
  | "PendingReport"
  | "Completed";
