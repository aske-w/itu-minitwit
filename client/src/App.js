import './App.css';
import { Routes, Route, Link, useNavigate } from "react-router-dom";
import Index from "./pages/Index"
import Public from './pages/Public';
import Signin from "./pages/Signin"
import Signup from "./pages/Signup"
import { useSelector, useDispatch } from 'react-redux'
import { logout } from './reducers/auth'

function App() {
	const dispatch = useDispatch()
	const auth = useSelector(state => state.auth)
	const navigate = useNavigate()

	const handleSignout = () => {
		dispatch(logout())

		navigate("/public")
	}

	return (
		<div className="page">
			<h1>MiniTwit</h1>

			<div className="navigation">
				{ auth.isLoggedIn ?
					<>
						<Link to="/">my timeline</Link> |
						<Link to="/public">public timeline</Link> |
						<button type="button" onClick={handleSignout}>sign out [{auth.user.username}]</button>
					</>
				:
					<>
						<Link to="/public">Public Timeline</Link> |
						<Link to="/signup">Sign Up</Link> |
						<Link to="/signin">Sign In</Link>
					</>
				}
			</div>

			<div className="body">
				<Routes>
					<Route path='/' element={<Index/>} />
					<Route path='/public' element={<Public/>} />
					{/* <Route exact path ="/">
						{!auth.isLoggedIn ? <Navigate to="/public"/> : <Index />}
					</Route> */}
					<Route path='/signup' element={<Signup/>} />
					<Route path='/signin' element={<Signin/>} />
				</Routes>
			</div>

			<div className="footer">MiniTwit &mdash; An Iris Application</div>
		</div>
	);
}

export default App;
