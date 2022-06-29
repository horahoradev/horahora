import { PrivacyPolicy } from "./privacy-policy";
import { TermsOfService } from "./terms-of-service";

export function Footer(props) {
  const { userData, dataless } = props;

  return (
    <nav className="h-8 w-full bg-white dark:bg-gray-900 shadow">
      <div className="flex justify-around p-auto">
        <PrivacyPolicy />
        <TermsOfService />
      </div>
    </nav>
  );
}
