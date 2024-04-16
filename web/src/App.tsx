import { RouterProvider, createBrowserRouter } from 'react-router-dom';
import './App.css';
import HomePage from './pages/home-page';

const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
  },
]);

export default function App() {
  return (
    <RouterProvider router={router} />
  );
}
