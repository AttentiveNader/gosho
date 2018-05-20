function Loading(bool){
	var loading = document.getElementById("loading")
	if (bool){
		loading.style.display = "block"
		return
	}	
	loading.style.display = "none"
	
}
function closeErr(id){
	document.getElementById(id).style.display = "none"
}
function PostUrl(){
	console.log("hello")
	var err = document.getElementById("respErr");
	var shortenerr = document.getElementById("shortenErr");
	var	el = document.getElementById("url");
	var valid = validate({website: el.value}, {website: {url: true}});
	if (valid == undefined) {
		err.style.display = "none"

		var resp = document.getElementById("response");
		shortenerr.innerHTML = "";
		resp.innerHTML="";
		Loading(true);
		   var body = {
                 longurl : el.value
              }

        axios({
              method: 'post',
              url: '/post',
              data: body
            })
            .then(function (res) {
            	Loading()
                if (res.data.id) {
	  			resp.innerHTML += Pragraphize(Anchorize(res.data.shorturl,res.data.shorturl));
	  			}
	  			console.log(res)
              })
              .catch(function (error) {
              	err.style.display = "block"
              	shortenerr.innerHTML += Pragraphize(error.Message)
          		Loading()
            });
	

	}else{
		err.style.display ="block";
		shortenerr.innerHTML += Pragraphize("Url not valid");
	}
}
function Anchorize(str,href){
	return '<a href="'+href+'">'+str+'</a>'
}
function Pragraphize(str){
	return "<p>"+str+"</p>"
}