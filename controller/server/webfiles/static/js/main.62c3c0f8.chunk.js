(this["webpackJsonpencodarr-react-frontend"]=this["webpackJsonpencodarr-react-frontend"]||[]).push([[0],{36:function(e,t,a){},62:function(e,t,a){},81:function(e,t,a){},83:function(e,t,a){},86:function(e,t,a){},88:function(e,t,a){},89:function(e,t,a){},90:function(e,t,a){},91:function(e,t,a){"use strict";a.r(t);a(57);var s=a(1),i=a.n(s),r=a(27),n=a.n(r),c=(a(62),a(11)),o=a(12),l=a(14),d=a(13),h=a(31),j=a(21),b=a(10),u=a.n(b),p=a(19),x=a(29),O=a(24),v=a(8),m=a(56),f=a(46),g=(a(81),a(0));function k(e){return Object(g.jsx)("img",{className:"queue-icon",src:e.location,alt:e.alt,height:"20px",title:e.title})}var _=a.p+"static/media/headphones.d0211b61.svg";function y(){return Object(g.jsx)(k,{location:_,alt:"Headphones",title:"An additional stereo audio track will be created"})}var C=a.p+"static/media/play_button.a97890f5.svg";function w(){return Object(g.jsx)(k,{location:C,alt:"Play Button",title:"File will be encoded to HEVC"})}a(83),a(36);var S=a.p+"static/media/Info-I.4eaee4f4.svg",N=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).timerID=void 0,s.state={jobs:[],waitingOnServer:!0,showModal:!1,waitingRunnersText:""},s.timerID=setTimeout((function(){}),Number.POSITIVE_INFINITY),clearInterval(s.timerID),s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){var e=this;this.tick(),this.timerID=setInterval((function(){return e.tick()}),2e3)}},{key:"componentWillUnmount",value:function(){clearInterval(this.timerID)}},{key:"tick",value:function(){var e=this;u.a.get("/api/web/v1/running").then((function(t){var a=t.data.jobs;void 0!==a?(a.sort((function(e,t){return parseFloat(e.status.percentage)>parseFloat(t.status.percentage)?-1:1})),e.setState({jobs:a,waitingOnServer:!1})):console.error("Response from /api/web/v1/running returned undefined for data.jobs")})).catch((function(e){console.error("Request to /api/web/v1/running failed with error: ".concat(e))})),u.a.get("/api/web/v1/waitingrunners").then((function(t){if(0===t.data.Runners.length)e.setState({waitingRunnersText:"No waiting runners"});else{var a=t.data.Runners.toString();1!==t.data.Runners.length&&(a=a.slice(1)),e.setState({waitingRunnersText:a})}})).catch((function(e){console.error("Request to /api/web/v1/waitingrunners failed with error: ".concat(e))}))}},{key:"render",value:function(){var e=this,t=function(){return e.setState({showModal:!1})},a=this.state.jobs.map((function(e){return Object(g.jsx)(I,{fps:e.status.fps,uuid:e.job.uuid,filename:e.job.path,progress:e.status.percentage,runnerName:e.runner_name,stageValue:e.status.stage,jobElapsedTime:e.status.job_elapsed_time,stageElapsedTime:e.status.stage_elapsed_time,stageEstimatedTimeRemaining:e.status.stage_estimated_time_remaining,videoOperation:e.job.parameters.hevc,audioOperation:e.job.parameters.stereo},e.job.uuid)}));return Object(g.jsxs)("div",{children:[Object(g.jsx)("img",{className:"info-i",src:S,alt:"",height:"20px",onClick:function(){return e.setState({showModal:!0})}}),0!==a.length?a:Object(g.jsx)("h5",{className:"text-center",children:"No running jobs"}),Object(g.jsxs)(v.a,{show:this.state.showModal,onHide:t,children:[Object(g.jsx)(v.a.Header,{closeButton:!0,children:Object(g.jsx)(v.a.Title,{children:"Waiting Runners"})}),Object(g.jsx)(v.a.Body,{children:this.state.waitingRunnersText}),Object(g.jsx)(v.a.Footer,{children:Object(g.jsx)(p.a,{variant:"secondary",onClick:t,children:"Close"})})]})]})}}]),a}(i.a.Component);function I(e){return Object(g.jsxs)("div",{children:[Object(g.jsxs)(x.a,{style:{padding:"1rem"},children:[Object(g.jsxs)(x.a.Header,{className:"text-center",children:[Object(g.jsxs)("div",{className:"file-image-container",children:[Object(g.jsx)("h5",{children:e.filename}),Object(g.jsxs)("div",{className:"svg-flex-container",children:[e.videoOperation?Object(g.jsx)(w,{}):null,e.audioOperation?Object(g.jsx)(y,{}):null]})]}),Object(g.jsxs)("h6",{children:["Stage: ",e.stageValue]}),Object(g.jsxs)("h6",{children:["Runner: ",e.runnerName]})]}),Object(g.jsx)(m.a,{className:"progress-bar-style",animated:!0,now:parseFloat(e.progress),label:"".concat(e.progress,"%")}),Object(g.jsxs)(f.a,{children:[Object(g.jsx)(O.a,{children:Object(g.jsx)("h6",{className:"text-right",children:"Job Elapsed Time:"})}),Object(g.jsx)(O.a,{children:Object(g.jsx)("p",{children:e.jobElapsedTime})}),Object(g.jsx)(O.a,{children:Object(g.jsx)("h6",{className:"text-right",children:"FPS:"})}),Object(g.jsx)(O.a,{children:Object(g.jsx)("p",{children:e.fps})})]}),Object(g.jsxs)(f.a,{children:[Object(g.jsx)(O.a,{children:Object(g.jsx)("h6",{className:"text-right",children:"Stage Elapsed Time:"})}),Object(g.jsx)(O.a,{children:Object(g.jsx)("p",{children:e.stageElapsedTime})}),Object(g.jsx)(O.a,{children:Object(g.jsx)("h6",{className:"text-right",children:"Stage Estimated Time Remaining:"})}),Object(g.jsx)(O.a,{children:Object(g.jsx)("p",{children:e.stageEstimatedTimeRemaining})})]})]}),Object(g.jsx)("div",{className:"smol-spacer"})]})}var T=a(28),H=a(15),V=a(3),E=a(40),P=(a(86),a.p+"static/media/addLibraryIcon.22f57e9e.svg"),R=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).timerID=void 0,s.state={libraries:[],waitingOnServer:!0,showCreateLibModal:!1},s.timerID=setTimeout((function(){}),Number.POSITIVE_INFINITY),clearInterval(s.timerID),s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){var e=this;this.tick(),this.timerID=setInterval((function(){return e.tick()}),2e3)}},{key:"componentWillUnmount",value:function(){clearInterval(this.timerID)}},{key:"tick",value:function(){var e=this;u.a.get("/api/web/v1/libraries").then((function(t){200===t.status&&e.setState({libraries:t.data.IDs})})).catch((function(e){console.error("Request to /api/web/v1/libraries failed with error: ".concat(e))}))}},{key:"render",value:function(){var e=this,t=this.state.libraries.map((function(e){return Object(g.jsxs)("div",{children:[Object(g.jsx)(F,{id:e}),Object(g.jsx)("div",{className:"smol-spacer"})]},e)}));return Object(g.jsxs)(g.Fragment,{children:[Object(g.jsx)("img",{className:"add-lib-ico",src:P,alt:"",height:"20px",onClick:function(){e.setState({showCreateLibModal:!0})}}),Object(g.jsx)(L,{show:this.state.showCreateLibModal,closeHandler:function(){e.setState({showCreateLibModal:!1})}}),Object(g.jsx)("div",{className:"smol-spacer"}),t]})}}]),a}(i.a.Component),F=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={folder:"",priority:"",fs_check_interval:"",path_masks:"",queue:[],target_video_codec:"HEVC",create_stereo_audio:!0,skip_hdr:!0,showEditModal:!1,showQueueModal:!1},s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){this.getLibraryData()}},{key:"getLibraryData",value:function(){var e=this;u.a.get("/api/web/v1/library/".concat(this.props.id)).then((function(t){e.setState({folder:t.data.folder,priority:t.data.priority,fs_check_interval:t.data.fs_check_interval,path_masks:t.data.path_masks.join(","),queue:t.data.queue.Items,target_video_codec:t.data.pipeline.target_video_codec,create_stereo_audio:t.data.pipeline.create_stereo_audio,skip_hdr:t.data.pipeline.skip_hdr})})).catch((function(t){console.error("Request to /api/web/v1/library/".concat(e.props.id," failed with error: ").concat(t))}))}},{key:"render",value:function(){var e=this;return Object(g.jsxs)(g.Fragment,{children:[Object(g.jsxs)(x.a,{children:[Object(g.jsx)(x.a.Header,{className:"text-center",children:Object(g.jsx)("h5",{children:this.state.folder})}),Object(g.jsxs)("p",{className:"text-center",children:["Priority: ",this.state.priority]}),Object(g.jsxs)("p",{className:"text-center",children:["File System Check Interval: ",this.state.fs_check_interval]}),Object(g.jsxs)("p",{className:"text-center",children:["Target Video Codec: ",this.state.target_video_codec]}),Object(g.jsxs)("p",{className:"text-center",children:["Create Stereo Audio Track: ",this.state.create_stereo_audio?"True":"False"]}),Object(g.jsxs)("p",{className:"text-center",children:["Skip HDR Files: ",this.state.skip_hdr?"True":"False"]}),0!==this.state.path_masks.length?Object(g.jsxs)("p",{className:"text-center",children:["Path Masks: ",this.state.path_masks]}):null,Object(g.jsx)(p.a,{variant:"secondary",onClick:function(){e.setState({showQueueModal:!0})},children:"Queue"}),Object(g.jsx)(p.a,{variant:"primary",onClick:function(){e.setState({showEditModal:!0})},children:"Edit"})]}),this.state.showEditModal?Object(g.jsx)(D,{show:!0,closeHandler:function(){e.setState({showEditModal:!1}),e.getLibraryData()},id:this.props.id,folder:this.state.folder,priority:this.state.priority,fs_check_interval:this.state.fs_check_interval,path_masks:this.state.path_masks,target_video_codec:this.state.target_video_codec,create_stereo_audio:this.state.create_stereo_audio,skip_hdr:this.state.skip_hdr}):null,this.state.showQueueModal?Object(g.jsx)(M,{show:!0,closeHandler:function(){e.setState({showQueueModal:!1})},queue:this.state.queue}):null]})}}]),a}(i.a.Component),L=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={folder:"",priority:"",fs_check_interval:"",path_masks:"",target_video_codec:"HEVC",create_stereo_audio:!0,skip_hdr:!0},s.submitLib=s.submitLib.bind(Object(T.a)(s)),s}return Object(o.a)(a,[{key:"submitLib",value:function(){var e=this,t={folder:this.state.folder,priority:parseInt(this.state.priority),fs_check_interval:this.state.fs_check_interval,path_masks:this.state.path_masks.split(","),pipeline:{target_video_codec:this.state.target_video_codec,create_stereo_audio:this.state.create_stereo_audio,skip_hdr:this.state.skip_hdr}};u.a.post("/api/web/v1/library/new",t).then((function(){e.props.closeHandler()})).catch((function(e){console.error("/api/web/v1/library/new failed with error: ".concat(e))}))}},{key:"render",value:function(){var e=this;return Object(g.jsx)("div",{children:Object(g.jsxs)(v.a,{show:this.props.show,onHide:this.props.closeHandler,children:[Object(g.jsx)(v.a.Header,{closeButton:!0,children:Object(g.jsx)(v.a.Title,{children:"Create New Library"})}),Object(g.jsxs)(v.a.Body,{children:[Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Folder"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"/home/user/lib1","aria-label":"folder","aria-describedby":"basic-addon1",onChange:function(t){e.setState({folder:t.target.value})},value:this.state.folder})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Priority"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"0","aria-label":"priority","aria-describedby":"basic-addon1",onChange:function(t){e.setState({priority:t.target.value})},value:this.state.priority})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"File System Check Interval"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"0h0m0s","aria-label":"fs_check_interval","aria-describedby":"basic-addon1",onChange:function(t){e.setState({fs_check_interval:t.target.value})},value:this.state.fs_check_interval})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Target Video Codec"})}),Object(g.jsxs)(H.a,{className:"dark-text-input no-box-shadow",as:"select",custom:!0,onChange:function(t){e.setState({target_video_codec:t.target.value})},value:this.state.target_video_codec,children:[Object(g.jsx)("option",{value:"HEVC",children:"H.265 (HEVC)"}),Object(g.jsx)("option",{value:"AVC",children:"H.264 (AVC)"}),Object(g.jsx)("option",{value:"VP9",children:"VP9"})]})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Create Stereo Audio Track"})}),Object(g.jsx)(V.a.Checkbox,{"aria-label":"Create Stereo Audio Track Checkbox",onChange:function(t){e.setState({create_stereo_audio:t.target.checked})},checked:this.state.create_stereo_audio})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Skip HDR"})}),Object(g.jsx)(V.a.Checkbox,{"aria-label":"Skip HDR Checkbox",onChange:function(t){e.setState({skip_hdr:t.target.checked})},checked:this.state.skip_hdr})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Path Masks"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"Plex Versions,private,.m4a","aria-label":"path_masks","aria-describedby":"basic-addon1",onChange:function(t){e.setState({path_masks:t.target.value})},value:this.state.path_masks})]})]}),Object(g.jsxs)(v.a.Footer,{children:[Object(g.jsx)(p.a,{variant:"secondary",onClick:this.props.closeHandler,children:"Close"}),Object(g.jsx)(p.a,{variant:"primary",onClick:this.submitLib,children:"Create"})]})]})})}}]),a}(i.a.Component),D=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={folder:e.folder,priority:e.priority,fs_check_interval:e.fs_check_interval,path_masks:e.path_masks,target_video_codec:e.target_video_codec,create_stereo_audio:e.create_stereo_audio,skip_hdr:e.skip_hdr},s.putChanges=s.putChanges.bind(Object(T.a)(s)),s.deleteLibrary=s.deleteLibrary.bind(Object(T.a)(s)),s}return Object(o.a)(a,[{key:"putChanges",value:function(){var e=this,t={folder:this.state.folder,priority:parseInt(this.state.priority),fs_check_interval:this.state.fs_check_interval,path_masks:this.state.path_masks.split(","),pipeline:{target_video_codec:this.state.target_video_codec,create_stereo_audio:this.state.create_stereo_audio,skip_hdr:this.state.skip_hdr}};u.a.put("/api/web/v1/library/".concat(this.props.id),t).then((function(){e.props.closeHandler()})).catch((function(t){console.error("/api/web/v1/library/".concat(e.props.id," failed with error: ").concat(t))}))}},{key:"deleteLibrary",value:function(){var e=this;u.a.delete("/api/web/v1/library/".concat(this.props.id)).then((function(){e.props.closeHandler()})).catch((function(t){console.error("/api/web/v1/library/".concat(e.props.id," failed with error: ").concat(t))}))}},{key:"render",value:function(){var e=this;return Object(g.jsx)("div",{children:Object(g.jsxs)(v.a,{show:this.props.show,onHide:this.props.closeHandler,children:[Object(g.jsx)(v.a.Header,{closeButton:!0,children:Object(g.jsx)(v.a.Title,{children:"Edit Library"})}),Object(g.jsxs)(v.a.Body,{children:[Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Folder"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"/home/user/lib1","aria-label":"folder","aria-describedby":"basic-addon1",onChange:function(t){e.setState({folder:t.target.value})},value:this.state.folder})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Priority"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"0","aria-label":"priority","aria-describedby":"basic-addon1",onChange:function(t){e.setState({priority:t.target.value})},value:this.state.priority})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"File System Check Interval"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"0h0m0s","aria-label":"fs_check_interval","aria-describedby":"basic-addon1",onChange:function(t){e.setState({fs_check_interval:t.target.value})},value:this.state.fs_check_interval})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Target Video Codec"})}),Object(g.jsxs)(H.a,{className:"dark-text-input no-box-shadow",as:"select",custom:!0,onChange:function(t){e.setState({target_video_codec:t.target.value})},value:this.state.target_video_codec,children:[Object(g.jsx)("option",{value:"HEVC",children:"H.265 (HEVC)"}),Object(g.jsx)("option",{value:"AVC",children:"H.264 (AVC)"}),Object(g.jsx)("option",{value:"VP9",children:"VP9"})]})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Create Stereo Audio Track"})}),Object(g.jsx)(V.a.Checkbox,{"aria-label":"Create Stereo Audio Track Checkbox",onChange:function(t){e.setState({create_stereo_audio:t.target.checked})},checked:this.state.create_stereo_audio})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Skip HDR"})}),Object(g.jsx)(V.a.Checkbox,{"aria-label":"Skip HDR Checkbox",onChange:function(t){e.setState({skip_hdr:t.target.checked})},checked:this.state.skip_hdr})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Path Masks"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"Plex Versions,private,.m4a","aria-label":"path_masks","aria-describedby":"basic-addon1",onChange:function(t){e.setState({path_masks:t.target.value})},value:this.state.path_masks})]})]}),Object(g.jsxs)(v.a.Footer,{children:[Object(g.jsx)(p.a,{className:"delete-button",variant:"danger",onClick:this.deleteLibrary,children:"Delete"}),Object(g.jsx)(p.a,{variant:"secondary",onClick:this.props.closeHandler,children:"Close"}),Object(g.jsx)(p.a,{variant:"primary",onClick:this.putChanges,children:"Update"})]})]})})}}]),a}(i.a.Component),M=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(){return Object(c.a)(this,a),t.apply(this,arguments)}return Object(o.a)(a,[{key:"render",value:function(){var e=this.props.queue;null===e&&(e=[]);var t=e.map((function(e,t){return Object(g.jsx)(A,{index:t+1,path:e.path,videoOperation:e.parameters.encode,audioOperation:e.parameters.stereo},e.uuid)}));return Object(g.jsx)("div",{children:Object(g.jsxs)(v.a,{show:this.props.show,onHide:this.props.closeHandler,size:"lg",children:[Object(g.jsx)(v.a.Header,{closeButton:!0,children:Object(g.jsx)(v.a.Title,{children:"Queue"})}),Object(g.jsx)(v.a.Body,{children:Object(g.jsxs)(E.a,{children:[Object(g.jsx)("thead",{children:Object(g.jsxs)("tr",{children:[Object(g.jsx)("th",{scope:"col",children:"#"}),Object(g.jsx)("th",{scope:"col",children:"File"})]})}),Object(g.jsx)("tbody",{children:t})]})}),Object(g.jsx)(v.a.Footer,{children:Object(g.jsx)(p.a,{variant:"secondary",onClick:this.props.closeHandler,children:"Close"})})]})})}}]),a}(i.a.Component);function A(e){return Object(g.jsxs)("tr",{children:[Object(g.jsx)("th",{scope:"row",children:e.index}),Object(g.jsx)("td",{children:e.path}),Object(g.jsx)("td",{children:Object(g.jsxs)("div",{className:"queue-icon-container",children:[e.videoOperation?Object(g.jsx)("span",{className:"play-button-image",children:Object(g.jsx)(w,{})}):null,e.audioOperation?Object(g.jsx)(y,{}):null]})})]})}var q=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).timerID=void 0,s.state={jobs:[],waitingOnServer:!0},s.timerID=setTimeout((function(){}),Number.POSITIVE_INFINITY),clearInterval(s.timerID),s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){var e=this;this.tick(),this.timerID=setInterval((function(){return e.tick()}),2e3)}},{key:"componentWillUnmount",value:function(){clearInterval(this.timerID)}},{key:"tick",value:function(){var e=this;u.a.get("/api/web/v1/history").then((function(t){var a=t.data.history;void 0!==a?(a.reverse(),e.setState({jobs:a,waitingOnServer:!1})):console.error("Response from /api/web/v1/history returned undefined for data.history")})).catch((function(e){console.error("Request to /api/web/v1/history failed with error: ".concat(e))}))}},{key:"render",value:function(){var e=this.state.jobs.map((function(e,t){return Object(g.jsx)(B,{datetime:e.datetime_completed,file:e.file},t)})),t=Object(g.jsxs)("tr",{children:[Object(g.jsx)("th",{scope:"row",children:"-"}),Object(g.jsx)("td",{children:"Waiting on server"})]});return Object(g.jsx)(x.a,{children:Object(g.jsxs)(E.a,{hover:!0,size:"sm",children:[Object(g.jsx)("thead",{children:Object(g.jsxs)("tr",{children:[Object(g.jsx)("th",{scope:"col",children:"Time Completed"}),Object(g.jsx)("th",{scope:"col",children:"File"})]})}),Object(g.jsx)("tbody",{children:this.state.waitingOnServer?t:e})]})})}}]),a}(i.a.Component);function B(e){return Object(g.jsxs)("tr",{children:[Object(g.jsx)("td",{children:e.datetime}),Object(g.jsx)("td",{children:e.file})]})}a(88);var K=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={controller_version:"Could not contact a ring",web_api_versions:[],runner_api_versions:[]},s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){var e=this;u.a.get("/api").then((function(t){e.setState({web_api_versions:t.data.web.versions,runner_api_versions:t.data.runner.versions})})).catch((function(e){console.error("Request to /api failed with error: ".concat(e))})),u.a.get("/version").then((function(t){e.setState({controller_version:t.data})})).catch((function(e){console.error("Request to /api failed with error: ".concat(e))}))}},{key:"render",value:function(){return Object(g.jsxs)(g.Fragment,{children:[Object(g.jsx)("h5",{children:"About Encodarr"}),Object(g.jsxs)("p",{children:[Object(g.jsx)("b",{children:"License:"})," This project is licensed under the Mozilla Public License 2.0 a copy of which can be found ",Object(g.jsx)("a",{href:"https://github.comhttps://github.com/BrenekH/encodarr/blob/master/LICENSE",target:"_blank",rel:"noreferrer",children:"here"}),")"]}),Object(g.jsxs)("p",{children:[Object(g.jsx)("b",{children:"Controller Version:"})," ",this.state.controller_version]}),Object(g.jsx)("p",{className:"list-title",children:Object(g.jsx)("b",{children:"Supported API Versions:"})}),Object(g.jsxs)("ul",{className:"api-list",children:[Object(g.jsxs)("li",{children:[Object(g.jsx)("b",{children:"Web:"})," ",this.state.web_api_versions.join(", ")]}),Object(g.jsxs)("li",{children:[Object(g.jsx)("b",{children:"Runner:"})," ",this.state.runner_api_versions.join(", ")]})]}),Object(g.jsxs)("p",{children:[Object(g.jsx)("b",{children:"GitHub Repository:"})," ",Object(g.jsx)("a",{href:"https://github.com/BrenekH/encodarr",target:"_blank",rel:"noreferrer",children:"https://github.com/BrenekH/encodarr"})]})]})}}]),a}(i.a.Component),W=(a(89),function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(e){var s;return Object(c.a)(this,a),(s=t.call(this,e)).state={inputValues:{smallerFilesCheck:!1,healthCheckInterval:"",unresponsiveRunnerTimeout:"",logVerbosity:""},showSavedIndicator:!1},s.handleClick=s.handleClick.bind(Object(T.a)(s)),s}return Object(o.a)(a,[{key:"componentDidMount",value:function(){this.updateSettings()}},{key:"createChangeHandler",value:function(e){var t=this,a=arguments.length>1&&void 0!==arguments[1]&&arguments[1],s=function(s){var i=Object.assign({},t.state.inputValues);i[e]=a?s.target.checked:s.target.value,t.setState({inputValues:i})};return s.bind(this),s}},{key:"handleClick",value:function(){var e=this;u.a.put("/api/web/v1/settings",{SmallerFiles:this.state.inputValues.smallerFilesCheck,HealthCheckInterval:this.state.inputValues.healthCheckInterval,HealthCheckTimeout:this.state.inputValues.unresponsiveRunnerTimeout,LogVerbosity:this.state.inputValues.logVerbosity}).then((function(t){t.status>=200&&t.status<=299?e.setState({showSavedIndicator:!0}):console.error(t),e.updateSettings()}))}},{key:"updateSettings",value:function(){var e=this;u.a.get("/api/web/v1/settings").then((function(t){e.setState({inputValues:{smallerFilesCheck:t.data.SmallerFiles,healthCheckInterval:t.data.HealthCheckInterval,unresponsiveRunnerTimeout:t.data.HealthCheckTimeout,logVerbosity:t.data.LogVerbosity}})}))}},{key:"render",value:function(){var e=this,t=this.state.showSavedIndicator?Object(g.jsx)(Q,{}):null;return this.state.showSavedIndicator&&setTimeout((function(){e.setState({showSavedIndicator:!1})}),5e3),Object(g.jsxs)(g.Fragment,{children:[Object(g.jsxs)("div",{children:[Object(g.jsx)("h5",{children:"General"}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsxs)(V.a.Prepend,{children:[Object(g.jsx)(V.a.Text,{children:"Prefer Smaller Files"}),Object(g.jsx)(V.a.Checkbox,{"aria-label":"Prefer Smaller Files Checkbox",onChange:this.createChangeHandler("smallerFilesCheck",!0),checked:this.state.inputValues.smallerFilesCheck})]}),Object(g.jsx)(H.a,{className:"dark-text-input",value:"Only keep encoded files if they are smaller than the original.",disabled:!0,readOnly:!0})]}),Object(g.jsx)("div",{className:"spacer"}),Object(g.jsx)("h5",{children:"Runner Health"}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Runner Health Check Interval"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"0h0m0s","aria-label":"health-check-interval","aria-describedby":"basic-addon1",onChange:this.createChangeHandler("healthCheckInterval"),value:this.state.inputValues.healthCheckInterval})]}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Unresponsive Runner Timeout"})}),Object(g.jsx)(H.a,{className:"dark-text-input",placeholder:"0h0m0s","aria-label":"unresponsive-runner-timeout","aria-describedby":"basic-addon1",onChange:this.createChangeHandler("unresponsiveRunnerTimeout"),value:this.state.inputValues.unresponsiveRunnerTimeout})]}),Object(g.jsx)("div",{className:"spacer"}),Object(g.jsx)("h5",{children:"Logging"}),Object(g.jsxs)(V.a,{className:"mb-3",children:[Object(g.jsx)(V.a.Prepend,{children:Object(g.jsx)(V.a.Text,{children:"Log Verbosity"})}),Object(g.jsxs)(H.a,{className:"dark-text-input no-box-shadow",as:"select",custom:!0,onChange:this.createChangeHandler("logVerbosity"),value:this.state.inputValues.logVerbosity,children:[Object(g.jsx)("option",{value:"TRACE",children:"Trace"}),Object(g.jsx)("option",{value:"DEBUG",children:"Debug"}),Object(g.jsx)("option",{value:"INFO",children:"Info"}),Object(g.jsx)("option",{value:"WARNING",children:"Warning"}),Object(g.jsx)("option",{value:"ERROR",children:"Error"}),Object(g.jsx)("option",{value:"CRITICAL",children:"Critical"})]})]}),Object(g.jsx)("div",{className:"smol-spacer"}),Object(g.jsx)(p.a,{variant:"light",onClick:this.handleClick,children:"Save"}),t]}),Object(g.jsx)("div",{className:"spacer"}),Object(g.jsx)(K,{})]})}}]),a}(i.a.Component));function Q(){return Object(g.jsx)("p",{className:"pop-in-out",style:{display:"inline"},children:"Saved!"})}a(90);var U=a.p+"static/media/Encodarr-Logo.b8d2df4b.svg";function G(){return Object(g.jsxs)("div",{className:"header-flex header-content text-center",children:[Object(g.jsx)("img",{src:U,alt:"",height:"60px",title:""}),Object(g.jsx)("h1",{children:"ncodarr"})]})}var z=function(e){Object(l.a)(a,e);var t=Object(d.a)(a);function a(){return Object(c.a)(this,a),t.apply(this,arguments)}return Object(o.a)(a,[{key:"handleSelect",value:function(e){switch(e){case"libraries":window.history.replaceState(void 0,"","/libraries"),document.title="Libraries - Encodarr";break;case"history":window.history.replaceState(void 0,"","/history"),document.title="History - Encodarr";break;case"settings":window.history.replaceState(void 0,"","/settings"),document.title="Settings - Encodarr";break;case"running":window.history.replaceState(void 0,"","/running"),document.title="Running - Encodarr"}}},{key:"render",value:function(){var e="running";switch(window.location.pathname){case"/libraries":e="libraries";break;case"/history":e="history";break;case"/settings":e="settings"}return Object(g.jsxs)("div",{className:"container",children:[Object(g.jsx)(G,{}),Object(g.jsxs)(h.a.Container,{id:"tab-nav",defaultActiveKey:e,transition:!1,onSelect:this.handleSelect,children:[Object(g.jsxs)(j.a,{fill:!0,variant:"pills",children:[Object(g.jsx)(j.a.Item,{children:Object(g.jsx)(j.a.Link,{eventKey:"running",children:"Running"})}),Object(g.jsx)(j.a.Item,{children:Object(g.jsx)(j.a.Link,{eventKey:"libraries",children:"Libraries"})}),Object(g.jsx)(j.a.Item,{children:Object(g.jsx)(j.a.Link,{eventKey:"history",children:"History"})}),Object(g.jsx)(j.a.Item,{children:Object(g.jsx)(j.a.Link,{eventKey:"settings",children:"Settings"})})]}),Object(g.jsx)("div",{className:"spacer"}),Object(g.jsxs)(h.a.Content,{children:[Object(g.jsx)(h.a.Pane,{eventKey:"running",mountOnEnter:!0,unmountOnExit:!0,children:Object(g.jsx)(N,{})}),Object(g.jsx)(h.a.Pane,{eventKey:"libraries",mountOnEnter:!0,unmountOnExit:!0,children:Object(g.jsx)(R,{})}),Object(g.jsx)(h.a.Pane,{eventKey:"history",mountOnEnter:!0,unmountOnExit:!0,children:Object(g.jsx)(q,{})}),Object(g.jsx)(h.a.Pane,{eventKey:"settings",mountOnEnter:!0,unmountOnExit:!0,children:Object(g.jsx)(W,{})})]})]}),Object(g.jsx)("div",{className:"smol-spacer"})]})}}]),a}(i.a.Component),J=function(e){e&&e instanceof Function&&a.e(3).then(a.bind(null,93)).then((function(t){var a=t.getCLS,s=t.getFID,i=t.getFCP,r=t.getLCP,n=t.getTTFB;a(e),s(e),i(e),r(e),n(e)}))};n.a.render(Object(g.jsx)(i.a.StrictMode,{children:Object(g.jsx)(z,{})}),document.getElementById("root")),J()}},[[91,1,2]]]);
//# sourceMappingURL=main.62c3c0f8.chunk.js.map