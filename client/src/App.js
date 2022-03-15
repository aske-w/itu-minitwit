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
		<div className="min-h-screen bg-gray-50">
			<div className="flex max-w-4xl mx-auto py-12">
				<div className="w-1/4 shrink-0">
					<nav className="space-y-1" aria-label="Sidebar">
						{ auth.isLoggedIn ?
							<>
								<Link
									to="/"
									className="flex items-center px-3 py-2 text-sm font-medium rounded-md"
								>My timeline</Link>

								<Link
									to="/public"
									className="flex items-center px-3 py-2 text-sm font-medium rounded-md"
								>Public timeline</Link>

								<button
									type="button"
									className='flex items-center px-3 py-2 text-sm font-medium rounded-md'
									onClick={handleSignout}
								>Sign out</button>
							</>
						:
							<>
								<Link
									to="/public"
									className="flex items-center px-3 py-2 text-sm font-medium rounded-md"
								>Public timeline</Link>

								<Link
									to="/signup"
									className="flex items-center px-3 py-2 text-sm font-medium rounded-md"
								>Sign up</Link>

								<Link
									to="/signin"
									className="flex items-center px-3 py-2 text-sm font-medium rounded-md"
								>Sign in</Link>
							</>
						}
					</nav>
				</div>

				<div className="w-3/4">
					<Routes>
						<Route path='/' element={<Index/>} />
						<Route path='/public' element={<Public/>} />
						<Route path='/signup' element={<Signup/>} />
						<Route path='/signin' element={<Signin/>} />
					</Routes>
				</div>
			</div>
		</div>
	);
}

export default App;
