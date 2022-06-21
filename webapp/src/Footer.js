import PrivacyPolicy from "./PrivacyPolicy";
import TermsOfService from "./TOS";

function Footer(props) {
  const { userData, dataless } = props;

  return (
    <nav className="h-8 w-full bg-white dark:bg-gray-900 shadow">
        <div className="flex justify-around p-auto">
            <div><PrivacyPolicy></PrivacyPolicy></div>
            <div><TermsOfService></TermsOfService></div>
        </div>
    </nav>
  );
}

export default Footer;
