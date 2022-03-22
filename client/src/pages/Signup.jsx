import React, { useState } from 'react'
import api from '../api'
import { login } from '../reducers/auth'
import { useNavigate } from "react-router-dom";
import { useDispatch } from 'react-redux'

const Signup = () => {
    const [form, setForm] = useState({
        username: "",
        email: "",
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

        api.post("/register", form)
            .then(() => {
                api.post("signin", {
                    username: form.username,
                    password: form.pwd,
                }).then(response => {
                    dispatch(login(response.data))

                    navigate("/")
                })
            })
            .catch(error => {
                if ([400, 422].includes(error.response.status)) {
                    setErrors(error.response.data.errors)
                }
            })
    }

    return (
        <div>
            <h2 className="text-2xl font-semibold mb-4">Sign up</h2>

            { errors.length > 0 &&
                <ul className="bg-red-100 px-3 py-2 my-4 rounded text-red-400 text-sm">
                    { errors.map((error) => <li key={error}>{error}</li>) }
                </ul>
            }

            <form className="mt-4" action="POST" onSubmit={handleSubmit}>
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
                                className="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
                                placeholder="Chose a username"
                                onChange={handleChange}
                                autoComplete="off"
                            />
                        </div>
                    </div>

                    <div>
                        <label htmlFor="email" className="block text-sm font-medium text-gray-700">
                            Email
                        </label>
                        <div className="mt-1">
                            <input
                                type="email"
                                name="email"
                                id="email"
                                className="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
                                placeholder="you@example.com"
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
                                className="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md"
                                placeholder="Enter a password"
                                onChange={handleChange}
                                autoComplete="off"
                            />
                        </div>
                    </div>
                </div>

                <button
                    type="submit"
                    className="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                    Sign up
                </button>
            </form>
        </div>
    )
}

export default Signup