import React, { Component } from 'react';

class QuizForm extends Component {

    state = {
      isAdded:false,
      selectedOption: 'one answer'
    }
  

    handleClick = () => {
      this.setState(
        state => ({
          ...state,
          isAdded: true
        })
      );
      console.log("question added")
    }

    change = (event) => {
      this.setState(
        state => ({
          ...state,
          selectedOption: event.target.value
        })
      )
      console.log(this.state.selectedOption)
    }
    render() {
      
      return (
        <div>
          <div>
            {this.state.isAdded ? 
           
            <div className="border border-purple-500 text-purple-800 text-xs flex flex-col justify-center items-center py-2 my-4">
                <label htmlFor="question">Choose type of the question</label>
                <select name="question" id="question" onChange={this.change}>
                  <option value="one answer">One answer</option>
                  <option value="multiple answers">Multiple answers</option>
                  <option value="true/false">True/False</option>
                  <option value="open">Open</option>
                </select> 

                <div className="flex flex-row items items-center pb-2">
                  <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                      Question 1
                  </label>
                  <input className=" h-6 text-gray-700 border border-purple-400 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="question" type="text" />
              </div>
              <div className="flex flex-row items items-center pb-2">  
                <label className="block uppercase tracking-wide text-gray-700 text-xs mb-2 px-4 pt-1" >
                      option 1
                  </label>
                  <input className="w-1/2 h-6 text-gray-700 border border-purple-400 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="option1" type="text" />
              </div>  
              <div className="flex flex-row items items-center pb-2">  
                <label className="block uppercase tracking-wide text-gray-700 text-xs mb-2 px-4 pt-1" >
                      option 2
                  </label>
                  <input className="w-1/2 h-6 text-gray-700 border border-purple-400 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="option2" type="text" />
              </div>  
              <div className="flex flex-row items items-center pb-2">  
                <label className="block uppercase tracking-wide text-gray-700 text-xs mb-2 px-4 pt-1 " >
                      option 3
                  </label>
                  <input className="w-1/2 h-6 text-gray-700 border border-purple-400 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="option3" type="text" />
              </div>  
              <input className="shadow bg-purple-800 hover:bg-purple-500 focus:shadow-outline focus:outline-none text-white py-1 px-2 rounded" type="button" value="Save"/>
              
            </div>

            : null}
          </div>
          <div className="flex justify-center">
            <button type='button' className=" w-3/4 border border-purple-500 border-dashed rounded px-4 py-5 mb-3 flex justify-center items items-center text-xl text-purple-700" onClick={this.handleClick}>
              + Add Question
            </button>

          </div>

        </div>
        
      );
    }
  }

  export default QuizForm;