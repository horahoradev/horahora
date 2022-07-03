import { UploadForm } from "#components/forms/upload";
import { Page } from "#components/page";

function UploadPage() {
  return (
    <Page>
      <div className="flex justify-center mx-4">
        <UploadForm/>
      </div>
    </Page>
  );
}

export default UploadPage;
