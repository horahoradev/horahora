
import Footer from "./Footer";
import { ThemeSwitcherProvider } from 'react-css-theme-switcher';


function App() {
  const themes = {
    light: '/antd.min.css',
    dark: '/antd.dark.min.css',
  };

  return (
    <ThemeSwitcherProvider defaultTheme="dark" themeMap={themes}>
    <div className=" bg-yellow-50 dark:bg-gray-900 min-h-screen font-sans-serif">
      
    </div>
    <Footer></Footer>
    </ThemeSwitcherProvider>
  );
}

export default App;
