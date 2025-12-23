export async function getReportSummary() {
  const res = await fetch("http://localhost:8080/report/summary", {
    cache: "no-store",
  });
  return res.json();
}

export async function getReports() {
  const res = await fetch("http://localhost:8080/reports", {
    cache: "no-store",
  });
  return res.json();
}

export async function getReportById(id: number) {
  const res = await fetch(`http://localhost:8080/reports/${id}`, {
    cache: "no-store",
  });
  return res.json();
}

//function update
export async function updateReportStatus(
  id: number,
  statusId: number
) {
  const res = await fetch(`http://localhost:8080/reports/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      report_status: statusId,
    }),
  });

  if (!res.ok) {
    throw new Error("Update failed");
  }

  return res.json();
}
