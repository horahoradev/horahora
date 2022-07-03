import type { FormEvent } from "react";

export function UploadForm() {
  async function handleFileUpload(event: FormEvent<HTMLFormElement>) {}

  return <form onSubmit={handleFileUpload}></form>;
}
