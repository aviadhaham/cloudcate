import logo from "@/assets/favicon.ico";

export default function Header() {
  return (
    <div className="flex items-center justify-center mt-7">
      <img
        className="mr-3.5"
        src={logo}
        style={{ width: "45px", height: "45px" }}
      ></img>
      <h1 className="text-4xl font-bold">AWS Resource Search</h1>
    </div>
  );
}
