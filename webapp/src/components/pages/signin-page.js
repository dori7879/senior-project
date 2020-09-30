import React from 'react';

const SignInPage = () => {
    return(
        <div className="min-h-screen flex items-center justify-center bg-purple-100 py-8 px-4 sm:px-6 lg:px-8">
            <div className="border border-purple-300">
                <div className="max-w-md w-64 mx-4">
                    <div>
                        <h2 className="mt-6 text-center text-2xl uppercase leading-7 font-extrabold text-purple-900">
                            Login
                        </h2>
                    </div>
                    <form className="mt-8" action="#" method="POST" >
                        <input type="hidden" name="remember" value="true" />
                        <div className="rounded shadow-sm">
                            <div>
                            <input aria-label="Email address" name="email" type="email" required className="appearance-none rounded-none relative block w-full px-3 py-2 border border-purple-200 placeholder-purple-400 text-gray-800 rounded focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5" placeholder="Email address"/>
                            </div>
                            <div className="mt-4">
                            <input aria-label="Password" name="password" type="password" required className="appearance-none rounded-none relative block w-full px-3 py-2 border border-purple-200 placeholder-purple-400 text-gray-800 rounded focus:outline-none focus:shadow-outline-blue focus:border-blue-300 focus:z-10 sm:text-sm sm:leading-5" placeholder="Password" />
                            </div>
                        </div>

                        <div className="mt-6">
                            <div className="flex justify-center items-center">
                                <a href="#" className="font-medium text-purple-900 hover:text-purple-700 focus:outline-none focus:underline transition ease-in-out duration-150">
                                    Forgot your password?
                                </a>
                            </div>
                        </div>
                        <div className="mt-6">
                            <button type="submit" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-600 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                                Sign in
                            </button>
                        </div>
                    </form>
                </div>    
            </div>
        </div>
    )
}

export default SignInPage;