import Footer from '../footer';
import Header from '../header';
import React from 'react';
import { render } from 'react-dom';

class StudetnHwPage extends React.Component{
    
    constructor(props) {
        super(props);
        this.onChangeFullName = this.onChangeFullName.bind(this);
            
        this.state = {
          fullName: "",
          answer: "",
          files: [],
          submitDate: new Date(),
          successful: false
        };          
    }

    onChangeFullName(e) {
        this.setState({
          fullName: e.target.value,
        });
    }   

    render(){
        
        return(
            <div>
                <Header />
                <div className="min-h-screen flex items-center flex-col justify-top bg-purple-100 py-2 px-4 sm:px-6 lg:px-8 ">              
                    <div className="flex justify-center flex-col items-center pb-4">
                        <div className="text-purple-900 font-bold text-xl pt-4">Homework</div>
                        <div className="flex flex-row items items-center pb-2">
                            <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                                Student name
                            </label>
                            <input onChange={this.onChangeFullName} className="text-gray-700 border border-purple-400  rounded py-1 text-xs px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="Enter your full name" />
                        </div>
                        <h3 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">Course Title</h3>  
                        <h2 className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">HW title</h2> 
                            
                                     
                    </div>
                </div>   
                <Footer />
            </div>
        )
    }
    }

export default StudetnHwPage;