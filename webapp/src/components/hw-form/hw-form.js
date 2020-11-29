import CKEditor from 'ckeditor4-react';
import DateTimePicker from 'react-datetime-picker';
import React from 'react';
import { Redirect } from 'react-router-dom';
import { connect } from "react-redux";
import { createHomework } from "../../actions/homework";

CKEDITOR.config.autoParagraph = false;

class HwForm extends React.Component {
    
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.onChangeFullName = this.onChangeFullName.bind(this);
        this.onChangeCourseTitle= this.onChangeCourseTitle.bind(this);
        this.onChangeTitle = this.onChangeTitle.bind(this);
        this.onChangeDescription = this.onChangeDescription.bind(this);
        this.onChangeFiles = this.onChangeFiles.bind(this);
        this.onChangeOpenDate = this.onChangeOpenDate.bind(this);
        this.onChangeCloseDate = this.onChangeCloseDate.bind(this);

    
        this.state = {
          fullName: "",
          courseTitle: "",
          title: "",
          description: "",
          files: [],
          openDate: new Date(),
          closeDate: new Date(),
          successful: false,
          isClicked: false
        };          
    }

    onChangeFullName(e) {
        this.setState({
          fullName: e.target.value,
        });
    }   
    onChangeCourseTitle(e) {
        this.setState({
          courseTitle: e.target.value,
        });
    }
    onChangeTitle(e) {
        this.setState({
          title: e.target.value,
        });
    }    
    onChangeDescription(e) {
        this.setState({
        description: e.editor.getData()
        });
    }
    onChangeFiles(e) {
        this.setState({
          files: e.target.value[0],
         });
    }

    onChangeOpenDate(date) {
        this.setState({
          openDate: date,
         });
        
    }

    onChangeCloseDate(date) {
        this.setState({
          closeDate: date,
         });
    }

    handleSubmit(e){
        e.preventDefault();
        this.setState({
            isClicked:true
        })
        const { dispatch} = this.props;
        dispatch(createHomework( this.state.courseTitle, this.state.title, this.state.description, this.state.files, this.state.openDate, this.state.closeDate, this.state.fullName))
            .then(() => {
                this.setState({
                successful: true,
                });
            })
            .catch(() => {
                this.setState({
                successful: false,
                });
            });
    }
    render(){
        if (this.state.isClicked) {
            return <Redirect to="/link"/>;
        }
        return(
            <form onSubmit={this.handleSubmit} className="w-3/4 border border-purple-300 rounded bg-purple-300 p-4">
                <div className="flex flex-row items items-center pb-2">
                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                        Full name*
                    </label>
                    <input  onChange={this.onChangeFullName} className="text-gray-700 border border-purple-400 rounded text-xs py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="Enter your full name" />
                </div>
                <div className="flex flex-row items items-center pb-2">
                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                        Course title*
                    </label>
                    <input onChange={this.onChangeCourseTitle} className="text-gray-700 border border-purple-400  rounded py-1 text-xs px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="Enter course title" />
                </div>
                <div className="flex flex-row items items-center pb-2">
                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                        Title*
                    </label>
                    <input  onChange={this.onChangeTitle} className="text-gray-700 border border-purple-400 text-xs rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="title" type="text" placeholder="Enter the homework title" />
                </div>
                <div className="flex flex-row pb-2">
                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                        Description
                    </label>
                    <CKEditor
                        data={this.state.description}
                        onChange={this.onChangeDescription}
                    />
                </div>
                <div className="flex flex-row pb-2 items items-center">
                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                        Attach files
                    </label>
                    <input onChange={this.onChangeFiles} className="border border-purple-400 text-xs text-gray-700 w-1/2 rounded py-1 px-2 leading-tight focus:outline-none focus:bg-white" id="files" type="file" multiple />
                </div>
                <div className="flex flex-row pb-2 items items-center">
                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1" >
                        Open date and time*
                    </label>
                    <DateTimePicker value={this.state.openDate} onChange={this.onChangeOpenDate} />
                </div>
                <div className="flex flex-row pb-3 items items-center">
                    <label className="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2 px-4 pt-1">
                        Close date and time*
                    </label>
                    <DateTimePicker value={this.state.closeDate} onChange={this.onChangeCloseDate} />
                </div>
                <div className="flex justify-center">
                     <button type="submit" className="mb-2 relative w-full flex justify-center py-1 px-2 border border-transparent text-sm leading-4 font-medium rounded-md text-purple-200 bg-purple-800 hover:bg-purple-500 focus:outline-none transition duration-150 ease-in-out">
                       Generate Link
                    </button>
                </div>
            </form>
        ) 
    }
    
}
function mapStateToProps(state) {
    const { isLoggedIn } = state.auth;
    return {
      isLoggedIn
    };
  }

export default connect(mapStateToProps)(HwForm);