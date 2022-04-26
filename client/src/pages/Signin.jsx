import React, { useState } from 'react'
import api from "../api"
import { login } from '../reducers/auth'
import { useNavigate } from "react-router-dom";
import { useDispatch } from 'react-redux'

const Signin = () => {
	const [form, setForm] = useState({
        username: "",
        pwd: "",
    })

	const [errors, setErrors] = useState([])

	const navigate = useNavigate()
	const dispatch = useDispatch()

    const handleChange = (event) => {
        setForm({...form, [event.target.name]: event.target.value})
    }

    function handleSubmit(e) {
        e.preventDefault()

		api.post("signin", form).then(response => {
			dispatch(login(response.data))

			navigate("/")
        }).catch(error => {
			if ([400, 422].includes(error.response.status)) {
				setErrors(error.response.data.errors)
			}
		})
    }

  	return (
		<div>
			<h2 className="mb-4 text-2xl font-semibold">Sign in</h2>

			{ errors.length > 0 &&
                <ul className="px-3 py-2 my-4 text-sm text-red-400 bg-red-100 rounded">
                    { errors.map((error) => <li key={error}>{error}</li>) }
                </ul>
            }

			<form className="mt-4" onSubmit={handleSubmit}>
			<div className="space-y-4">
                    <div>
                        <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                            Username
                        </label>
                        <div className="mt-1">
                            <input
                                type="text"
                                name="username"
                                id="username"
                                className="block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                placeholder="Chose a username"
                                onChange={handleChange}
                                autoComplete="off"
                            />
                        </div>
                    </div>

                    <div>
                        <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                            Password
                        </label>
                        <div className="mt-1">
                            <input
                                type="password"
                                name="pwd"
                                id="password"
                                className="block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                placeholder="Enter a password"
                                onChange={handleChange}
                                autoComplete="off"
                            />
                        </div>
                    </div>
                </div>

				<button
                    type="submit"
                    className="inline-flex items-center px-4 py-2 mt-4 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                    Sign in
                </button>
			</form>
		</div>
  	)
}

export default Signin