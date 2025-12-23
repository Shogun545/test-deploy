export async function uploadProfileImage(file: File): Promise<string> {
  const formData = new FormData();
  formData.append("file", file);
  formData.append("upload_preset", "profile_image");

  const res = await fetch(
    "https://api.cloudinary.com/v1_1/de6bzhgjg/image/upload",
    {
      method: "POST",
      body: formData,
    }
  );

  const data = await res.json();
  return data.secure_url;
}
