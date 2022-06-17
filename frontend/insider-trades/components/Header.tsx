import Navbar from "./Navbar";
import Search from "./Search";

export const Header = () => {
  return (
    <header className="mt-14 mb-10 mx-4">
      <div className="container header bg-white flex md:flex-row flex-col justify-between items-center overflow-hidden">
        <a href="/" className="mx-6 logo-font">Insider trades</a>
        <Navbar />

        <Search />
      </div>
    </header>
  );
};
