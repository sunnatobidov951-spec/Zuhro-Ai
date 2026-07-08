function addProduct() {
    let name = document.getElementById("name").value;
    let price = document.getElementById("price").value;
    
    if (name !== "" && price !== "") {
        let table = document.getElementById("productTable");
        let row = table.insertRow(-1);
        let cell1 = row.insertCell(0);
        let cell2 = row.insertCell(1);
        cell1.innerHTML = name;
        cell2.innerHTML = price + " $";
        
        // Очистка полей ввода
        document.getElementById("name").value = "";
        document.getElementById("price").value = "";
    } else {
        alert("Пожалуйста, заполните оба поля!");
    }
}

