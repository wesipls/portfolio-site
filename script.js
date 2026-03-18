function prepareAndPrint() {
    const originalTitle = document.title;
    // This sets the filename for the "Save as PDF" dialog
    document.title = "Kim_Valkama_Infrastructure_and_Platform_Engineer_CV"; 
    
    window.print();
    
    // Restore the original browser tab title
    document.title = originalTitle;
}
