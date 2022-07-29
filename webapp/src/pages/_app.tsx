import "antd/dist/antd.css";
import "#styles/index.scss";

import Head from "next/head";
import type { AppProps } from "next/app";

import { Footer } from "#components/footer";
import { AccountProvider } from "#hooks";

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <meta charSet="utf-8" />
        <title>Horahora</title>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="theme-color" content="#000000" />
        <meta name="description" content="horahora" />
        <script
          id="theme-setup"
          async
          dangerouslySetInnerHTML={{
            __html: `
              const theme = document.cookie
                .split("; ")
                .find((row) => row.startsWith("theme="));
              const value = theme ? theme.split("=")[1] : "dark";
              document.documentElement.dataset.theme = value;`,
          }}
        ></script>
        <link rel="icon" href="/favicon.ico" />
        <link rel="apple-touch-icon" href="/logo192.png" />
        <link rel="manifest" href="/manifest.json" />
      </Head>
      <div className="text-base bg-yellow-50 dark:bg-gray-900 min-h-screen font-sans-serif">
        <AccountProvider>
          <Component {...pageProps} />
        </AccountProvider>
      </div>

      <Footer />
    </>
  );
}

export default MyApp;
