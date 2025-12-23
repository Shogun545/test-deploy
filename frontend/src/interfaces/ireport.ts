export interface Report {
  ID?: number;
  description?: string;
  createdAt: string;

  report_by?: number;

  report_status?: number;
  report_topic?: number;
  report_image?: number;
}