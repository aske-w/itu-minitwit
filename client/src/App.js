import './App.css';
import { Routes, Route, Link } from "react-router-dom";
import Index from './pages/Index';
import Signin from "./pages/Signin"
import Signup from "./pages/Signup"

function App() {
	return (
		<div className="page">
			<h1>MiniTwit</h1>

			<div className="navigation">
        		<Link to="/">Public Timeline</Link> |
        		<Link to="/signup">Sign Up</Link> |
        		<Link to="/signin">Sign In</Link>
			</div>

			<div className="body">
				<Routes>
					<Route path='/' element={<Index/>} />
					<Route path='/signup' element={<Signup/>} />
					<Route path='/signin' element={<Signin/>} />
				</Routes>
			</div>

			<div className="footer">MiniTwit &mdash; An Iris Application</div>
		</div>
	);
}

export default App;
