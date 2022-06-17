import SearchBar from "components/SearchBar";

function Navbar() {
  return (
    <div className="bg-slate-50 w-full shadow shadow-gray-300 max-w-screen-2xl mx-4 px-4 md:px-auto rounded-xl">
      <div className="flex flex-col items-center justify-center flex-wrap sm:flex-row">
        <div className="flex order-last pb-4 sm:pb-0 sm:order-first justify-between">
          <a className="text-indigo-500 px-3 hover:text-indigo-900" href="/">
            All
          </a>
          <a className="text-indigo-500 px-3 hover:text-indigo-900" href="/">
            Buy
          </a>
          <a className="text-indigo-500 px-3 hover:text-indigo-900" href="/">
            Sell
          </a>
        </div>

        <SearchBar />
      </div>
    </div>
  );
}

export default Navbar;
