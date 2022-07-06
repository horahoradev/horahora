import { LinkInternal } from "#components/links";

interface IFooterProps extends Record<string, unknown> {}

export function Footer({ userData, dataless }: IFooterProps) {
  return (
    <nav className="flex justify-around h-8 w-full bg-white dark:bg-gray-900 shadow">
      <LinkInternal href="/privacy-policy">Privacy Policy</LinkInternal>
      <LinkInternal href="/terms-of-service">Terms of Service</LinkInternal>
    </nav>
  );
}
